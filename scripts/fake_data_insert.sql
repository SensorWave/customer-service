-- =========================================
-- INSERT 10 COMPANIES
-- =========================================
INSERT INTO companies (name, email, phone, address) VALUES
('TechNova Solutions', 'contact@technova.com', '+33145010001', '12 Rue des Innovations, Paris'),
('GreenWorld Industries', 'info@greenworld.com', '+33155880002', '48 Avenue de la Nature, Lyon'),
('SmartLogiX', 'support@smartlogix.com', '+33155660003', '92 Boulevard du Transport, Lille'),
('BlueOcean Technologies', 'hello@blueocean.com', '+33144020004', '7 Quai des Vagues, Marseille'),
('NextGen Labs', 'contact@nextgenlabs.com', '+33133090005', '21 Rue du Futur, Toulouse'),
('UrbanConnect', 'support@urbanconnect.com', '+33122070006', '15 Avenue des Metropoles, Bordeaux'),
('CyberShield Security', 'info@cybershield.com', '+33177880007', '5 Rue de la Défense, Nanterre'),
('Solaris Energies', 'contact@solaris.com', '+33166770008', '33 Allée du Soleil, Montpellier'),
('AeroDynamics France', 'support@aerodynamics.com', '+33155990009', '60 Rue de l’Aviation, Nantes'),
('DataPulse Analytics', 'hello@datapulse.com', '+33133990010', '88 Boulevard des Données, Strasbourg');

-- =========================================
-- INSERT 50 USERS
-- =========================================

-- Company 1
INSERT INTO users (company_id, first_name, last_name, email, phone) VALUES
(1,'Alice','Martin','alice.martin@technova.com','+33600010001'),
(1,'Bruno','Duchamp','bruno.duchamp@technova.com','+33600010002'),
(1,'Camille','Robert','camille.robert@technova.com','+33600010003'),
(1,'Damien','Lefevre','damien.lefevre@technova.com','+33600010004'),
(1,'Eva','Morel','eva.morel@technova.com','+33600010005');

-- Company 2
INSERT INTO users (company_id, first_name, last_name, email, phone) VALUES
(2,'Claire','Dubois','claire.dubois@greenworld.com','+33600020001'),
(2,'David','Lambert','david.lambert@greenworld.com','+33600020002'),
(2,'Emma','Durand','emma.durand@greenworld.com','+33600020003'),
(2,'Florian','Guillot','florian.guillot@greenworld.com','+33600020004'),
(2,'Hélène','Garnier','helene.garnier@greenworld.com','+33600020005');

-- Company 3
INSERT INTO users (company_id, first_name, last_name, email, phone) VALUES
(3,'François','Morel','francois.morel@smartlogix.com','+33600030001'),
(3,'Isabelle','Jean','isabelle.jean@smartlogix.com','+33600030002'),
(3,'Julien','Bernard','julien.bernard@smartlogix.com','+33600030003'),
(3,'Karim','Said','karim.said@smartlogix.com','+33600030004'),
(3,'Laura','Masson','laura.masson@smartlogix.com','+33600030005');

-- Company 4
INSERT INTO users (company_id, first_name, last_name, email, phone) VALUES
(4,'Marc','Giraud','marc.giraud@blueocean.com','+33600040001'),
(4,'Nina','Perrot','nina.perrot@blueocean.com','+33600040002'),
(4,'Olivier','Roux','olivier.roux@blueocean.com','+33600040003'),
(4,'Pauline','Blanc','pauline.blanc@blueocean.com','+33600040004'),
(4,'Quentin','Paris','quentin.paris@blueocean.com','+33600040005');

-- Company 5
INSERT INTO users (company_id, first_name, last_name, email, phone) VALUES
(5,'Romain','Collet','romain.collet@nextgenlabs.com','+33600050001'),
(5,'Sarah','Dupuy','sarah.dupuy@nextgenlabs.com','+33600050002'),
(5,'Thomas','Legrand','thomas.legrand@nextgenlabs.com','+33600050003'),
(5,'Ugo','Barbier','ugo.barbier@nextgenlabs.com','+33600050004'),
(5,'Valérie','Picard','valerie.picard@nextgenlabs.com','+33600050005');

-- Company 6
INSERT INTO users (company_id, first_name, last_name, email, phone) VALUES
(6,'William','Lucas','william.lucas@urbanconnect.com','+33600060001'),
(6,'Xavier','Fabre','xavier.fabre@urbanconnect.com','+33600060002'),
(6,'Yasmine','Chevalier','yasmine.chevalier@urbanconnect.com','+33600060003'),
(6,'Zoé','Besson','zoe.besson@urbanconnect.com','+33600060004'),
(6,'Adrien','Hardy','adrien.hardy@urbanconnect.com','+33600060005');

-- Company 7
INSERT INTO users (company_id, first_name, last_name, email, phone) VALUES
(7,'Amandine','Royer','amandine.royer@cybershield.com','+33600070001'),
(7,'Bastien','Hoareau','bastien.hoareau@cybershield.com','+33600070002'),
(7,'Cindy','Guérin','cindy.guerin@cybershield.com','+33600070003'),
(7,'Dorian','Pascal','dorian.pascal@cybershield.com','+33600070004'),
(7,'Elodie','Adam','elodie.adam@cybershield.com','+33600070005');

-- Company 8
INSERT INTO users (company_id, first_name, last_name, email, phone) VALUES
(8,'Fabrice','Diallo','fabrice.diallo@solaris.com','+33600080001'),
(8,'Gina','Rolland','gina.rolland@solaris.com','+33600080002'),
(8,'Hugo','Meunier','hugo.meunier@solaris.com','+33600080003'),
(8,'Inès','Boulanger','ines.boulanger@solaris.com','+33600080004'),
(8,'Jean','Colin','jean.colin@solaris.com','+33600080005');

-- Company 9
INSERT INTO users (company_id, first_name, last_name, email, phone) VALUES
(9,'Katia','Maurel','katia.maurel@aerodynamics.com','+33600090001'),
(9,'Luc','Girard','luc.girard@aerodynamics.com','+33600090002'),
(9,'Manon','Prévost','manon.prevost@aerodynamics.com','+33600090003'),
(9,'Nicolas','Fontaine','nicolas.fontaine@aerodynamics.com','+33600090004'),
(9,'Oceane','Leclerc','oceane.leclerc@aerodynamics.com','+33600090005');

-- Company 10
INSERT INTO users (company_id, first_name, last_name, email, phone) VALUES
(10,'Pascal','Briand','pascal.briand@datapulse.com','+33600100001'),
(10,'Rita','Fournier','rita.fournier@datapulse.com','+33600100002'),
(10,'Samir','Diallo','samir.diallo@datapulse.com','+33600100003'),
(10,'Tania','Maury','tania.maury@datapulse.com','+33600100004'),
(10,'Yohan','Chartier','yohan.chartier@datapulse.com','+33600100005');
