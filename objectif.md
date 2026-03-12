# Contexte & décision d’architecture (Keycloak + Customer Service)

## Contexte
Nous avons actuellement :
- **Keycloak** déjà opérationnel côté front (Vue) pour l’authentification (login/mot de passe) et l’émission de JWT.
- Un **Customer Service (Go + Postgres)** qui gère le modèle métier : **companies** et logique associée.
- Un problème initial de **dédoublement** : Keycloak stocke des utilisateurs “identité”, et le Customer Service stocke aussi des utilisateurs.

Objectif : **un seul service de connexion** (Keycloak), tout en gardant une gestion métier “user ↔ company” dans le Customer Service, sans dupliquer l’identité.

## Décision
1) **Keycloak est la source de vérité pour l’identité et l’authentification**
    - login / mot de passe / MFA / reset password / sessions
    - identifiants d’utilisateur Keycloak (UUID) exposés dans les tokens (claim `sub`)

2) **Customer Service est la source de vérité pour le modèle métier**
    - appartenance à une company
    - rôle métier au sein de la company (ex: `company_owner`, `company_admin`, `member`)
    - autorisations métier (qui peut gérer les membres, etc.)

3) **Aucune authentification n’est faite contre la DB du Customer Service**
    - Keycloak ne “tape” pas dans la DB du Customer Service pour login/register.
    - Le backend valide le JWT Keycloak, lit `sub`, puis résout le contexte métier via sa DB.

## Modèle de données retenu
Dans le Customer Service, on ne stocke pas l’identité (email/prénom/nom/mot de passe) de façon obligatoire.
On stocke uniquement une table de liaison (ex. `company_members` ou `user_company`) :

- `keycloak_user_id` (string/UUID) = `sub` du JWT Keycloak
- `company_id`
- `role` (métier)

> Option : enrichir l’affichage (email/nom) à la demande via Keycloak Admin API, sans persister ces champs.

## Flux d’exécution (runtime)
### Auth / appels API
1. L’utilisateur se connecte sur **Keycloak** (front Vue via `keycloak-js`).
2. Le front appelle l’API avec `Authorization: Bearer <access_token>`.
3. Le backend :
    - valide le JWT (issuer + JWKS + exp)
    - extrait `sub`
    - résout `company_id` et `role` via la table `company_members`

### Endpoint minimal requis
- `GET /me` (protégé) :
    - lit `sub`
    - retourne `{ company_id, role }`
    - retourne `404 { code: "NOT_LINKED" }` si aucune liaison n’existe

## Gestion des membres (chef de company)
But à terme : un `company_owner` peut ajouter/supprimer des membres de sa company.

Décision :
- Le front “CRUD users” (chef) appelle **le Customer Service**
- Le Customer Service applique les règles métier (owner/admin seulement)
- Le Customer Service modifie la DB de liaisons (`company_members`)
- La création/suppression du compte Keycloak se fait **via le backend** (jamais depuis le front)

### Phasage (état actuel)
- Actuellement : les comptes Keycloak sont créés manuellement par un admin Keycloak.
- Le backend doit donc fournir un endpoint “bootstrap” interne pour créer la liaison DB :
    - `POST /admin/memberships/link { keycloak_user_id, company_id, role }`

### Phasage (évolution)
- Plus tard : le backend créera les comptes Keycloak depuis l’UI chef via Keycloak Admin API
    - service account (`backend-admin`) + `client_credentials`
    - création user + required actions (verify email / set password)
    - puis insertion de la liaison `company_members`

## Non-objectifs (pour éviter la complexité)
- Ne pas implémenter de “User Federation/SPI” Keycloak pour authentifier contre la DB du Customer Service.
- Ne pas exposer de credentials ou d’accès Admin Keycloak au front.
- Ne pas dupliquer systématiquement l’identité (email/nom) en DB métier.

## Conséquences
- Il existe deux “représentations” d’un user :
    - **Keycloak** (identité/sécurité)
    - **Customer Service** (métier : company/role via `keycloak_user_id`)
- Ce n’est pas un doublon d’auth : **un seul login** (Keycloak).
- Les autorisations “chef gère ses membres” sont contrôlées par le backend via DB métier.
