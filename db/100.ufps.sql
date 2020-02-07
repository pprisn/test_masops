
DROP TABLE IF EXISTS ufps;
CREATE TABLE ufps (
  id		VARCHAR(32)	PRIMARY KEY COMMENT 'Ид УФПС',
  name			VARCHAR(250)	NOT NULL COMMENT 'Название'
) COMMENT 'УФПС';


ALTER TABLE nsis ADD (
  ufpsid   VARCHAR(32)	COMMENT 'Ид УФПС'
);

INSERT INTO ufps (id,name) values ('C00', '');
INSERT INTO ufps (id,name) values ('R00', '');
INSERT INTO ufps (id,name) values ('R01', '');
INSERT INTO ufps (id,name) values ('R02', '');
INSERT INTO ufps (id,name) values ('R03', '');
INSERT INTO ufps (id,name) values ('R04', '');
INSERT INTO ufps (id,name) values ('R05', '');
INSERT INTO ufps (id,name) values ('R06', '');
INSERT INTO ufps (id,name) values ('R07', '');
INSERT INTO ufps (id,name) values ('R08', '');
INSERT INTO ufps (id,name) values ('R09', '');
INSERT INTO ufps (id,name) values ('R10', '');
INSERT INTO ufps (id,name) values ('R11', '');
INSERT INTO ufps (id,name) values ('R12', '');
INSERT INTO ufps (id,name) values ('R13', '');
INSERT INTO ufps (id,name) values ('R14', '');
INSERT INTO ufps (id,name) values ('R15', '');
INSERT INTO ufps (id,name) values ('R16', '');
INSERT INTO ufps (id,name) values ('R17', '');
INSERT INTO ufps (id,name) values ('R18', '');
INSERT INTO ufps (id,name) values ('R19', '');
INSERT INTO ufps (id,name) values ('R21', '');
INSERT INTO ufps (id,name) values ('R22', '');
INSERT INTO ufps (id,name) values ('R23', '');
INSERT INTO ufps (id,name) values ('R24', '');
INSERT INTO ufps (id,name) values ('R25', '');
INSERT INTO ufps (id,name) values ('R26', '');
INSERT INTO ufps (id,name) values ('R27', '');
INSERT INTO ufps (id,name) values ('R28', '');
INSERT INTO ufps (id,name) values ('R29', '');
INSERT INTO ufps (id,name) values ('R30', '');
INSERT INTO ufps (id,name) values ('R31', '');
INSERT INTO ufps (id,name) values ('R32', '');
INSERT INTO ufps (id,name) values ('R33', '');
INSERT INTO ufps (id,name) values ('R34', '');
INSERT INTO ufps (id,name) values ('R35', '');
INSERT INTO ufps (id,name) values ('R36', '');
INSERT INTO ufps (id,name) values ('R37', '');
INSERT INTO ufps (id,name) values ('R38', '');
INSERT INTO ufps (id,name) values ('R39', '');
INSERT INTO ufps (id,name) values ('R40', '');
INSERT INTO ufps (id,name) values ('R41', '');
INSERT INTO ufps (id,name) values ('R42', '');
INSERT INTO ufps (id,name) values ('R43', '');
INSERT INTO ufps (id,name) values ('R44', '');
INSERT INTO ufps (id,name) values ('R45', '');
INSERT INTO ufps (id,name) values ('R46', '');
INSERT INTO ufps (id,name) values ('R48', '');
INSERT INTO ufps (id,name) values ('R49', '');
INSERT INTO ufps (id,name) values ('R50', '');
INSERT INTO ufps (id,name) values ('R51', '');
INSERT INTO ufps (id,name) values ('R52', '');
INSERT INTO ufps (id,name) values ('R53', '');
INSERT INTO ufps (id,name) values ('R54', '');
INSERT INTO ufps (id,name) values ('R55', '');
INSERT INTO ufps (id,name) values ('R56', '');
INSERT INTO ufps (id,name) values ('R57', '');
INSERT INTO ufps (id,name) values ('R58', '');
INSERT INTO ufps (id,name) values ('R59', '');
INSERT INTO ufps (id,name) values ('R60', '');
INSERT INTO ufps (id,name) values ('R61', '');
INSERT INTO ufps (id,name) values ('R62', '');
INSERT INTO ufps (id,name) values ('R63', '');
INSERT INTO ufps (id,name) values ('R64', '');
INSERT INTO ufps (id,name) values ('R65', '');
INSERT INTO ufps (id,name) values ('R67', '');
INSERT INTO ufps (id,name) values ('R68', '');
INSERT INTO ufps (id,name) values ('R69', '');
INSERT INTO ufps (id,name) values ('R70', '');
INSERT INTO ufps (id,name) values ('R71', '');
INSERT INTO ufps (id,name) values ('R72', '');
INSERT INTO ufps (id,name) values ('R73', '');
INSERT INTO ufps (id,name) values ('R74', '');
INSERT INTO ufps (id,name) values ('R75', '');
INSERT INTO ufps (id,name) values ('R76', '');
INSERT INTO ufps (id,name) values ('R77', '');
INSERT INTO ufps (id,name) values ('R78', '');
INSERT INTO ufps (id,name) values ('R79', '');
INSERT INTO ufps (id,name) values ('R83', '');
INSERT INTO ufps (id,name) values ('R86', '');
INSERT INTO ufps (id,name) values ('R87', '');
INSERT INTO ufps (id,name) values ('R89', '');
INSERT INTO ufps (id,name) values ('R95', '');
INSERT INTO ufps (id,name) values ('R96', '');
UPDATE nsis SET ufpsid ="R48";
