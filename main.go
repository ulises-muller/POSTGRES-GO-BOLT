package main

import (
	"database/sql"
	"fmt"
	"log"
	_ "github.com/lib/pq"
	bolt "go.etcd.io/bbolt"
	"encoding/json"
	"strconv"
)

func main() {
	db1, err := sql.Open("postgres", "user=postgres host=localhost dbname=postgres sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}
	defer db1.Close()
	db2, err := sql.Open("postgres", "user=postgres host=localhost dbname=centro_medico sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}
	defer db2.Close()
	var opcion int
	for {
		fmt.Printf("\nQue quieres hacer ? ingresa la tecla indicada\n-Eliminar base de datos : 1\n-Crear base de datos : 2\n-Crear tablas : 3\n-Agregar keys : 4\n-Eliminar keys : 5\n-Insertar datos : 6\n-Borrar datos : 7\n-Stored Procedures y Triggers : 8\n-BoltDB : 9\n-Salir : 0\n")
		fmt.Scan(&opcion)
		if 1 == opcion {
			dropDatabase(db1)
		}
		if 2 == opcion {
			createDatabase(db1)
		}
		if 3 == opcion {
			crearTablas(db2)
		}
		if 4 == opcion {
			agregarKeys(db2)
		}
		if 5 == opcion {
			borrarKeys(db2)
		}
		if 6 == opcion {
			insertarDatos(db2)
		}
		if 7 == opcion {
			eliminarDatos(db2)
		}
		if 8 == opcion {
			triggersMenu(db2)
		}
		if 9 == opcion {
			boltStart(db2)
		}
		if 0 == opcion {
			fmt.Printf("Saliste del programa..\n")
			break
		}
	}
}
func triggersMenu(db *sql.DB){
	var opcion int
	for{
		fmt.Printf("\n> Stored Procedures y Triggers\n-Creacion : 1\n-Prueba : 2\n-Volver : 0\n")
		fmt.Scan(&opcion)
		if 1 == opcion {
			triggersCreate(db)
		}
		if 2 == opcion {
			triggersRun(db)
		}
		if 0 == opcion {
			fmt.Printf("Volviendo..")
			break
		}
	}
}
func triggersCreate(db *sql.DB){
	var opcion int
	for{
		fmt.Printf("\nMenu > Stored Procedures y Triggers > Menu de creacion\n-Generacion de turnos disponibles : 1\n-Reserva de turno : 2\n-Cancelacion de turnos : 3\n-Atencion de turnos : 4\n-Liquidacion para obras sociales : 5\n-Envio de emails a pacientes : 6\n-Volver : 0\n")
		fmt.Scan(&opcion)
		if 1 == opcion {
			generarTurnos(db)
		}
		if 2 == opcion {
			reservarFuncion(db)
		}
		if 3 == opcion {
			anularFuncion(db)
		}
		if 4 == opcion {
			atencionFuncion(db)
		}
		if 5 == opcion {
			liqObraSocialFuncion(db)
		}
		if 6 == opcion {
			envioEmail(db)
		}
		if 0 == opcion {
			fmt.Printf("Volviendo..")
			break
		}
	}
}
func triggersRun(db *sql.DB){
	var opcion int
	for{
		fmt.Printf("\nMenu > Stored Procedures y Triggers > Menu de pruebas\n-Generacion de turnos disponibles : 1\n-Reserva de turno : 2\n-Cancelacion de turnos : 3\n-Atencion de turnos : 4\n-Liquidacion para obras sociales : 5\n-Volver : 0\n")
		fmt.Scan(&opcion)
		if 1 == opcion {
			crearTurnos(db)
		}
		if 2 == opcion {
			reservarTurnos(db)
		}
		if 3 == opcion {
			anularTurnos(db)
		}
		if 4 == opcion {
			atencionTurnos(db)
		}
		if 5 == opcion {
			generarLiqObraSocial(db)
		}
		if 0 == opcion {
			fmt.Printf("Volviendo..")
			break
		}
	}
}
func dropDatabase(db *sql.DB) {
	_, err := db.Exec(`drop database centro_medico`)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Se elimino la bases de datos\n")
}
func createDatabase(db *sql.DB) {
	_, err := db.Exec(`create database centro_medico`)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Se creo la base de datos centro_medico\n")
}
func crearTablas(db *sql.DB) {
	_, err := db.Exec(`create table paciente(
                        nro_paciente int,
                        nombre text,
                        apellido text,
                        dni_paciente int,
                        f_nac date,
                        nro_obra_social int,
                        nro_afiliade int,
                        domicilio text,
                        telefono char(12),
                        email text
                    );
                    create table medique(
                        dni_medique int,
                        nombre text,
                        apellido text,
                        especialidad varchar(64),
                        monto_consulta_privada decimal(12, 2),
                        telefono char(12)
                    );create table consultorio(
                        nro_consultorio int,
                        nombre text,
                        domicilio text,
                        codigo_postal char(8),
                        telefono char(12)
                    );

                    create table agenda(
                        dni_medique int,
                        dia int,
                        nro_consultorio int,
                        hora_desde time,
                        hora_hasta time,
                        duracion_turno interval
                    );

                    create table turno(
                        nro_turno int,
                        fecha timestamp,
                        nro_consultorio int,
                        dni_medique int,
                        nro_paciente int,
                        nro_obra_social_consulta int,
                        nro_afiliade_consulta int,
                        monto_paciente decimal(12, 2),
                        monto_obra_social decimal(12, 2),
                        f_reserva timestamp,
                        estado char(10)
                    );

                    create table reprogramacion(
                        nro_turno int,
                        nombre_paciente text,
                        apellido_paciente text,
                        telefono_paciente char(12),
                        email_paciente text,
                        nombre_medique text,
                        apellido_medique text,
                        estado char(12)
                    );

                    create table error(
                        nro_error serial,
                        f_turno timestamp,
                        nro_consultorio int,
                        dni_medique int,
                        nro_paciente int,
                        operacion char(12),
                        f_error timestamp,
                        motivo varchar(64)
                    );

                    create table cobertura(
                        dni_medique int,
                        nro_obra_social int,
                        monto_paciente decimal(12, 2),
                        monto_obra_social decimal(12, 2)
                    );

                    create table obra_social (
                        nro_obra_social int,
                        nombre text,
                        contacto_nombre text,
                        contacto_apellido text,
                        contacto_telefono char(12),
                        contacto_email text
                    );

                    create table liquidacion_cabecera(
                        nro_liquidacion serial,
                        nro_obra_social int,
                        desde date,
                        hasta date,
                        total decimal(15, 2)
                    );

                    create table liquidacion_detalle(
                        nro_liquidacion int,
                        nro_linea serial,
                        f_atencion date,
                        nro_afiliade int,
                        dni_paciente int,
                        nombre_paciente text,
                        apellido_paciente text,
                        dni_medique int,
                        nombre_medique text,
                        apellido_medique text,
                        especialidad varchar(64),
                        monto decimal(12, 2)
                    );

                    create table envio_email(
                        nro_email serial,
                        f_generacion timestamp,
                        email_paciente text,
                        asunto text,
                        cuerpo text,
                        f_envio timestamp,
                        estado char(10)
                    );

                    create table solicitud_reservas(
                        nro_orden int,
                        nro_paciente int,
                        dni_medique int,
                        fecha timestamp
                    );`)
	fmt.Printf("Se agregaron las tablas\n")
	if err != nil {
		log.Fatal(err)
	}
}

func agregarKeys(db *sql.DB) {
	_, err := db.Exec(`alter table paciente add constraint pk_paciente primary key(nro_paciente);
						alter table medique add constraint pk_medique primary key(dni_medique);
						alter table consultorio add constraint pk_consultorio primary key(nro_consultorio);
						alter table turno add constraint pk_turno primary key(nro_turno);
						alter table error add constraint pk_error primary key(nro_error);
						alter table obra_social add constraint pk_obra_social primary key(nro_obra_social);
						alter table liquidacion_cabecera add constraint pk_liquidacion_cabecera primary key(nro_liquidacion);
						alter table liquidacion_detalle add constraint pk_liquidacion_detalle primary key(nro_linea);
						alter table envio_email add constraint pk_envio_email primary key(nro_email);
						alter table agenda add constraint fk_agenda_dni_medique foreign key(dni_medique) references medique(dni_medique);
						alter table agenda add constraint fk_agenda_nro_consultorio foreign key(nro_consultorio) references consultorio(nro_consultorio);
						alter table turno add constraint fk_turno_dni_medique foreign key(dni_medique) references medique(dni_medique);
						alter table turno add constraint fk_turno_nro_paciente foreign key(nro_paciente) references paciente(nro_paciente);
						alter table turno add constraint fk_turno_nro_consultorio foreign key(nro_consultorio) references consultorio(nro_consultorio);
						alter table reprogramacion add constraint fk_reprogramacion_nro_turno foreign key(nro_turno) references turno(nro_turno);
						alter table error add constraint fk_error_dni_medique foreign key(dni_medique) references medique(dni_medique);
						alter table error add constraint fk_error_nro_paciente foreign key(nro_paciente) references paciente(nro_paciente);
						alter table error add constraint fk_error_nro_consultorio foreign key(nro_consultorio) references consultorio(nro_consultorio);
						alter table cobertura add constraint fk_cobertura_dni_medique foreign key(dni_medique) references medique(dni_medique);
						alter table cobertura add constraint fk_cobertura_nro_obra_social foreign key(nro_obra_social) references obra_social(nro_obra_social);
						alter table liquidacion_cabecera add constraint fk_liquidacion_c_nro_obra_social foreign key(nro_obra_social) references obra_social(nro_obra_social);
						alter table liquidacion_detalle add constraint fk_liquidacion_d_liquidacion_c foreign key(nro_liquidacion) references liquidacion_cabecera(nro_liquidacion);
						alter table solicitud_reservas add constraint fk_solicitud_reservas_nro_paciente foreign key(nro_paciente) references paciente(nro_paciente);
						alter table solicitud_reservas add constraint fk_solicitud_reservas_dni_medique foreign key(dni_medique) references medique(dni_medique);
	`)
	fmt.Printf("Se insertaron las KEYS\n")
	if err != nil {
		log.Fatal(err)
	}
}

func borrarKeys(db *sql.DB) {
	_, err := db.Exec(`alter table agenda drop constraint fk_agenda_dni_medique; 
						alter table agenda drop constraint fk_agenda_nro_consultorio; 
						alter table turno drop constraint fk_turno_dni_medique; 
						alter table turno drop constraint fk_turno_nro_paciente; 
						alter table turno drop constraint fk_turno_nro_consultorio; 
						alter table reprogramacion drop constraint fk_reprogramacion_nro_turno; 
						alter table error drop constraint fk_error_dni_medique; 
						alter table error drop constraint fk_error_nro_paciente; 
						alter table error drop constraint fk_error_nro_consultorio; 
						alter table cobertura drop constraint fk_cobertura_dni_medique; 
						alter table cobertura drop constraint fk_cobertura_nro_obra_social; 
						alter table liquidacion_cabecera drop constraint fk_liquidacion_c_nro_obra_social; 
						alter table liquidacion_detalle drop constraint fk_liquidacion_d_liquidacion_c; 
						alter table solicitud_reservas drop constraint fk_solicitud_reservas_nro_paciente;
						alter table solicitud_reservas drop constraint fk_solicitud_reservas_dni_medique;
						alter table paciente drop constraint pk_paciente; 
						alter table medique drop constraint pk_medique;
						alter table consultorio drop constraint pk_consultorio;  
						alter table turno drop constraint pk_turno; 
						alter table error drop constraint pk_error; 
						alter table obra_social drop constraint pk_obra_social; 
						alter table liquidacion_cabecera drop constraint pk_liquidacion_cabecera; 
						alter table liquidacion_detalle drop constraint pk_liquidacion_detalle; 
						alter table envio_email drop constraint pk_envio_email; 
	`)
	fmt.Printf("Se borrarion las KEYS\n")
	if err != nil {
		log.Fatal(err)
	}
}

func insertarDatos(db *sql.DB) {
	result, err := db.Exec(`insert into paciente values (1, 'Juan', 'Perez', 12345678, '1978-05-08', 721, 523456, 'Suipacha 123', '+1153213421', 'juanperez1@gmail.com'),
                        	(2, 'Maria','Rodriguez', 23456789, '1980-06-09', 722, 234567, 'Av. Libertador 123', '+1153213422', 'mariarodriguez1@gmail.com'),
                        	(3, 'Pedro', 'Gomez', 34567890, '1982-07-10', 723, 345678, 'Calle 123', '+1153213423', 'pedrogomez1@gmail.com'),
                        	(4, 'Lucia', 'Fernandez', 45678901, '1984-08-11', 723, 456789, 'Calle 456', '+1153213424', 'luciafernandez1@gmail.com'),
                        	(5, 'Jorge', 'Gonzalez', 46789012, '1986-09-12', 722, 567890, 'Calle 789', '+1153213425', 'jorgegonzalez1@gmail.com'),
                        	(6, 'Ana', 'Martinez', 37890123, '1988-10-13', 722, 678901, 'Calle 012', '+1153213426', 'anamartinez1@gmail.com'),
                        	(7, 'Carlos', 'Sanchez', 28901234, '1990-11-14', 723, 789012, 'Calle 345', '+1153213427', 'carlossanchez1@gmail.com'),
                        	(8, 'Laura', 'Romero', 19012345, '1992-12-15', 721, 890123, 'Calle 678', '+1153213428', 'lauraromero1@gmail.com'),
                        	(9, 'Federico', 'Diaz', 20123456, '1994-01-16', 721, 901234, 'Calle 901', '+1153213429', 'federicodiaz1@gmail.com'),
                        	(10, 'Mariana', 'Castro', 32345670, '1996-02-17', 722, 123459, 'Calle 234', '+1153213430', 'marianacastro1@gmail.com'),
                        	(11, 'Roberto', 'Alvarez', 42345678, '1998-03-18', 721, 234569, 'Av. Libertador 456', '+1153213431', 'robertoalvarez1@gmail.com'),
                        	(12, 'Sofia', 'Acosta', 33456789, '2000-04-19', 722, 345677, 'Calle 789', '+1153213432', 'sofiaacosta1@gmail.com'),
                        	(13, 'Martin', 'Torres', 24567890, '2002-05-20', 723, 456784, 'Calle 012', '+1153213433', 'martintorres1@gmail.com'),
                        	(14, 'Valentina', 'Ruiz', 15678901, '2004-06-21', 721, 567895, 'Calle 345', '+1153213434', 'valentinaruiz1@gmail.com'),
                        	(15, 'Agustin', 'Sosa', 26789012, '2006-07-22', 722, 678905, 'Calle 678', '+1153213435', 'agustinsosa1@gmail.com'),
                        	(16, 'Camila', 'Castro', 37890123, '2008-08-23', 723, 789015, 'Calle 901', '+1153213436', 'camilacastro1@gmail.com'),
                        	(17, 'Lucas', 'Fernandez', 48901234, '2010-09-24', 721, 890125, 'Calle 234', '+1153213437', 'lucasfernandez1@gmail.com'),
                        	(18, 'Sofia', 'Gonzalez', 39012345, '2012-10-25', 722, 901235, 'Calle 567', '+1153213438', 'sofiagonzalez1@gmail.com'),
                        	(19, 'Luisa', 'Perez', 22345670, '2002-05-20', 723, 345666, 'Calle 345', '+1153213439', 'luisaperez1@gmail.com'),
                        	(20, 'Miguel', 'Rodriguez', 13456780, '2004-06-21', 721, 678747, 'Calle 678', '+1153213440', 'miguelrodriguez1@gmail.com');`)
	fmt.Printf("Se insertaron los pacientes\n")
	if err != nil {
		log.Fatal(err)
	}
	rowsAffected, err := result.RowsAffected()
    if err != nil {
        log.Fatal(err)
    }
    fmt.Printf("Se insertaron %d filas.\n", rowsAffected)
	result, err = db.Exec(`insert into medique values (44458951, 'Lara', 'Dolores', 'Traumatologo', 3000.50, '+1153223425'),
                        	(44458952, 'Juan', 'Cardozo', 'Cardiologo', 5000.00, '+1153233426'),
                        	(44458953, 'Sofia', 'Rodriguez', 'Pediatra', 2500.00, '+1153243427'),
                        	(44458954, 'Martin', 'Gonzalez', 'Dermatologo', 4000.00, '+1153253428'),
                        	(44458955, 'Ana', 'Fernandez', 'Ginecologo', 3500.00, '+1153216429'),
                        	(44458956, 'Diego', 'Maradona', 'Oftalmologo', 3200.00, '+1153273430'),
                        	(44458957, 'Valentina', 'Sanchez', 'Psiquiatra', 4500.00, '+1153813431'),
                        	(44458958, 'Lucas', 'Garcia', 'Neurologo', 3800.00, '+1153213492'),
                        	(44458959, 'Camila', 'Lopez', 'Oncologo', 4200.00, '+1153213733'),
                        	(44458960, 'Mateo', 'Diaz', 'Endocrinologo', 3600.00, '+1153213634'),
                        	(44458961, 'Agustina', 'Torres', 'Infectologo', 3300.00, '+1153289435'),
                        	(44458962, 'Cristina', 'Kirchner', 'Reumatologo', 4100.00, '+1153218836'),
                        	(44458963, 'Micaela', 'Castro', 'Nutricionista', 2800.00, '+1153671437'),
                        	(44458964, 'Tomas', 'Romero', 'Urologo', 3400.00, '+1153213434'),
                        	(44458965, 'Julieta', 'Gomez', 'Oncologo', 4400.00, '+1153211239'),
                        	(44458966, 'Luciana', 'Pereyra', 'Cardiologo', 4800.00, '+1153234540'),
                        	(44458967, 'Ignacio', 'Gutierrez', 'Traumatologo', 3100.00, '+1123163441'),
                        	(44458968, 'Florencia', 'Alvarez', 'Pediatra', 2900.00, '+1153211232'),
                        	(44458969, 'Santiago', 'Rojas', 'Dermatologo', 3700.00, '+1153243243'),
                        	(44458970, 'Lucia', 'Garcia', 'Oftalmologo', 3900.00, '+11532132344');`)
	fmt.Printf("Se insertaron los mediques\n")
	if err != nil {
		log.Fatal(err)
	}
	rowsAffected, err = result.RowsAffected()
    if err != nil {
        log.Fatal(err)
    }
    fmt.Printf("Se insertaron %d filas.\n", rowsAffected)
	result, err = db.Exec(`insert into consultorio values (1, 'Consultorio 1', 'Calle Consultorio 1', '1619', '1138624484'),
							(2, 'Consultorio 2', 'Calle Consultorio 2', '1620', '1157489652'),
							(3, 'Consultorio 3', 'Calle Consultorio 3', '1621', '1198654712'),
							(4, 'Consultorio 4', 'Calle Consultorio 4', '1622', '1185479623'),
							(5, 'Consultorio 5', 'Calle Consultorio 5', '1623', '1135987462'),
							(6, 'Consultorio 6', 'Calle Consultorio 6', '1624', '1155847999'),
							(7, 'Consultorio 7', 'Calle Consultorio 7', '1625', '1152366487');
	`)
	fmt.Printf("Se insertaron consultorios\n")
	if err != nil {
		log.Fatal(err)
	}
	rowsAffected, err = result.RowsAffected()
    if err != nil {
        log.Fatal(err)
    }
    fmt.Printf("Se insertaron %d filas.\n", rowsAffected)
	result, err = db.Exec(`insert into obra_social values (721, 'OSDE', 'Roberto', 'De Magallanes', 1154785299, 'contacto@osde.com.ar'),
							(722, 'Swiss Medical', 'Fernanda', 'Rodriguez', 1122447586, 'contacto@swissmedical.com.ar'),
							(723, 'Galeno', 'Federico', 'Campos', 1125334578, 'contacto@galenoargentina.com.ar');
	`)
	fmt.Printf("Se insertaron las obras sociales\n")
	if err != nil {
		log.Fatal(err)
	}
	rowsAffected, err = result.RowsAffected()
    if err != nil {
        log.Fatal(err)
    }
    fmt.Printf("Se insertaron %d filas.\n", rowsAffected)
	
	result, err = db.Exec(`insert into cobertura values (44458951, 721, 6000.00, 9000.00),
							(44458952, 722, 20000.00, 80000.00),
							(44458953, 723, 1000.00, 4000.00),
							(44458954, 723, 4500.00, 45500.00),
							(44458955, 722, 1500.00, 13500.00),
							(44458956, 722, 25000.00, 45000.00),
							(44458957, 723, 3000.00, 0.00),
							(44458958, 721, 12000.00, 78000.00),
							(44458959, 721, 20000.00, 80000.00),
							(44458960, 722, 4000.00, 56000.00),
							(44458961, 721, 0.00, 3300.00),
							(44458962, 722, 0.00, 12000.00),
							(44458963, 723, 5000.00, 0.00),
							(44458964, 721, 2500.00, 7500.00),
							(44458965, 722, 7500.00, 925000.00),
							(44458966, 723, 15000.00, 85000.00),
							(44458967, 722, 5000.00, 10000.00),
							(44458968, 722, 0.00, 5000.00),
							(44458969, 723, 4500.00, 5500.00),
							(44458970, 721, 70000.00, 0.00);
	`)
	fmt.Printf("Se insertaron las coberturas\n")
	if err != nil {
		log.Fatal(err)
	}
	rowsAffected, err = result.RowsAffected()
    if err != nil {
        log.Fatal(err)
    }
    fmt.Printf("Se insertaron %d filas.\n", rowsAffected)
	
	result, err = db.Exec(`insert into agenda values (44458951, 1, 1, '06:00', '12:00', '1 hour'),
							(44458951, 4, 2, '06:00', '12:00', '1 hour'),
							(44458952, 2, 3, '06:00', '12:00', '30 minutes'),
							(44458952, 5, 4, '06:00', '12:00', '30 minutes'),
							(44458953, 3, 5, '06:00', '12:00', '15 minutes'),
							(44458953, 6, 6, '06:00', '12:00', '15 minutes'),
							(44458954, 7, 7, '06:00', '12:00', '30 minutes'),
							(44458955, 5, 1, '12:00', '18:00', '1 hour'),
							(44458956, 6, 6, '12:00', '18:00', '1 hour'),
							(44458957, 3, 2, '12:00', '18:00', '30 minutes'),
							(44458958, 4, 3, '12:00', '18:00', '15 minutes'),
							(44458959, 1, 4, '12:00', '18:00', '40 minutes'),
							(44458960, 2, 5, '12:00', '18:00', '1 hour'),
							(44458961, 4, 1, '18:00', '00:00', '1 hour'),
							(44458962, 5, 2, '18:00', '00:00', '30 minutes'),
							(44458963, 2, 3, '18:00', '00:00', '15 minutes'),
							(44458964, 3, 4, '18:00', '00:00', '30 minutes'),
							(44458965, 1, 5, '18:00', '00:00', '40 minutes'),
							(44458966, 2, 1, '00:00', '06:00', '30 minutes'),
							(44458967, 3, 2, '00:00', '06:00', '1 hour'),
							(44458968, 4, 3, '00:00', '06:00', '1 hour'),
							(44458969, 5, 4, '00:00', '06:00', '30 minutes'),
							(44458970, 1, 5, '00:00', '06:00', '1 hour');
	`)
	fmt.Printf("Se insertaron los datos en la agenda\n")
	if err != nil {
		log.Fatal(err)
		}
	rowsAffected, err = result.RowsAffected()
    if err != nil {
        log.Fatal(err)
    }
    fmt.Printf("Se insertaron %d filas.\n", rowsAffected)
    
    result, err = db.Exec(`insert into solicitud_reservas values (1, 1, 44458951, '2023-11-22 08:00:00'),
							(2, 2, 44458952, '2023-11-16 11:00:00'),
							(3, 3, 44458953, '2023-11-01 07:15:00'),
							(4, 4, 44458954, '2023-12-14 06:30:00'),
							(5, 5, 44458955, '2023-12-26 12:00:00'),
							(6, 8, 44458959, '2023-12-29 16:00:00'),
							(7, 7, 44458957, '2023-12-17 15:00:00'),
							(8, 8, 44458958, '2023-11-04 17:45:00'),
							(9, 9, 44458959, '2023-12-22 14:40:00'),
							(10, 10, 44458960, '2023-11-09 14:00:00'),
							(11, 11, 44458961, '2024-01-19 18:00:00'),
							(12, 12, 44458962, '2024-11-28 19:30:00'),
							(13, 13, 44458963, '2024-12-17 23:45:00'),
							(14, 14, 44458964, '2024-12-16 20:00:00'),
							(15, 15, 44458965, '2024-01-30 21:20:00'),
							(16, 16, 44458966, '2023-11-02 02:30:00'),
							(17, 17, 44458967, '2023-12-25 04:00:00'),
							(18, 18, 44458968, '2023-11-25 05:00:00'),
							(19, 19, 44458969, '2023-12-05 04:30:00'),
							(20, 20, 44458970, '2023-11-08 00:00:00'),
							(21, 2, 44458952, '2023-11-02 06:00:00'),
							(22, 2, 44458955, '2023-12-12 17:00:00'),
							(23, 10, 44458952, '2023-11-02 08:00:00'),
							(24, 2, 44458952, '2023-11-09 08:00:00'),
							(25, 3, 44458953, '2023-11-24 12:00:00'),
							(26, 3, 44458953, '2023-11-24 11:00:00'),
							(27, 3, 44458953, '2023-11-24 10:00:00'),
							(28, 1, 44458951, '2023-11-01 08:00:00'),
							(29, 16, 44458953, '2023-11-24 11:00:00'),
							(30, 2, 44458952, '2023-11-26 08:00:00'),
							(31, 12, 44458952, '2023-11-26 06:00:00');
	`)
	fmt.Printf("Se insertaron las obras sociales\n")
	if err != nil {
		log.Fatal(err)
	}
	rowsAffected, err = result.RowsAffected()
    if err != nil {
        log.Fatal(err)
    }
    fmt.Printf("Se insertaron %d filas.\n", rowsAffected)
}
func eliminarDatos(db *sql.DB) {
	_ , err := db.Exec(`truncate table paciente cascade;
						truncate table medique cascade;
						truncate table agenda cascade;
						truncate table cobertura cascade;
						truncate table obra_social cascade;
						truncate table consultorio cascade;;
	`)
	fmt.Printf("Se eliminaron los datos\n")
	if err != nil {
		log.Fatal(err)
	}
}

func generarTurnos(db *sql.DB) {
	_, err := db.Exec(`create or replace function generar_turnos_disponibles(anio int, mes int)
returns boolean as $$
declare
	_mediquedni cursor for select distinct dni_medique from agenda order by dni_medique;
    _nroturno int = 0;
    _dnis record;
    _basura record;
    _horas record;
	_dias record;
	_dias2 timestamp;
    _fecha timestamp;
    _fecha2 timestamp;
    _fecha3 timestamp;
    _fecha4 timestamp;
    _fecha5 timestamp;
	_hora time;
begin
    select t.fecha into _basura from turno t
    where extract(year from t.fecha) = anio and extract(month from t.fecha) = mes;
    if found then
        return false;
    else
        if exists (select nro_turno from turno) then
            select max(nro_turno) into _nroturno from turno;
        end if;
        for _dnis in _mediquedni
        loop
                for _dias in (select dia from agenda a where a.dni_medique = _dnis.dni_medique)
            loop
                select distinct a.hora_desde, a.hora_hasta, a.duracion_turno into _horas from agenda a where a.dni_medique = _dnis.dni_medique;
                _fecha4 = to_timestamp(anio || '-' || mes || '-' || _dias.dia || ' ' || _horas.hora_desde, 'YYYY-MM-DD HH24:MI:SS');
                if mes != 12 then
                    _fecha5 = to_timestamp(anio || '-' || mes + 1 || '-' || _dias.dia || ' ' || _horas.hora_hasta, 'YYYY-MM-DD HH24:MI:SS');
                else
                    _fecha5 = to_timestamp(anio + 1 || '-' || 1 || '-' || _dias.dia || ' ' || _horas.hora_hasta, 'YYYY-MM-DD HH24:MI:SS');
                end if;
                for _dias2 in select generate_series(_fecha4::timestamp, _fecha5::timestamp, '7 days') as _dias2
                loop
                        _fecha = to_timestamp(anio || '-' || mes || '-' || extract(days from _dias2) || ' ' || _horas.hora_desde, 'YYYY-MM-DD HH24:MI:SS');
                        _fecha2 = to_timestamp(anio || '-' || mes || '-' || extract(days from _dias2) || ' ' || _horas.hora_hasta, 'YYYY-MM-DD HH24:MI:SS');
                        for _hora in select generate_series(_fecha::timestamp, _fecha2::timestamp, _horas.duracion_turno) as _hora
                        loop
                            _nroturno := _nroturno + 1;
                            _fecha3 = to_timestamp(anio || '-' || mes || '-' || extract(days from _dias2) || ' ' || _hora, 'YYYY-MM-DD HH24:MI:SS');
                            insert into turno(nro_turno, fecha, dni_medique, estado) values (_nroturno, _fecha3, _dnis.dni_medique, 'Disponible');
                        end loop;
                end loop;
            end loop;
        end loop;
        return true;
    end if;
end
$$ language plpgsql;
	`)
	if err != nil {
		log.Fatal(err)
	}
    fmt.Printf("CREATE FUNCTION\n")
}
func crearTurnos(db *sql.DB) {
	var opcion string
	var anio int
	var mes int
	fmt.Printf("Se crearon los turnos para noviembre del 2023\n")
	_, err := db.Query("select generar_turnos_disponibles(2023, 11)")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Quiere seguir creando turnos? y / n :  ")
	fmt.Scan(&opcion)
	if opcion == "y" {
		fmt.Printf("Ingresa año y mes separados con una , :  ")
		fmt.Scanf("%d, %d", &anio, &mes)
		rows, err := db.Query("select generar_turnos_disponibles($1, $2)", anio, mes)
		if err != nil {
			log.Fatal(err)
		}
		defer rows.Close()
		for rows.Next(){
			var resu bool
			err = rows.Scan(&resu)
			if err != nil {
				log.Fatal(err)
			}
			fmt.Printf("%b", resu)
		}
	} else { fmt.Printf("Volviendo..")}
	
}
func reservarFuncion(db *sql.DB) {
	_, err := db.Exec(`create or replace function reservar_turnos(_nro_paciente int, _dni_medique int, _fecha timestamp)
returns boolean as $$
declare
	_basura record;
begin
  select * from medique m into _basura where m.dni_medique in (_dni_medique);
  if found then
    select * from paciente p into _basura where p.nro_paciente in (_nro_paciente);
    if found then
      if exists (select * from paciente p, medique m, cobertura c where p.nro_obra_social = c.nro_obra_social and c.dni_medique = m.dni_medique and p.nro_paciente = _nro_paciente and m.dni_medique = _dni_medique) then
        if exists (select * from turno t where t.fecha in (_fecha) and t.estado = 'Disponible') then
            select t.nro_paciente into _basura from turno t
            group by t.nro_paciente having count (t.nro_paciente) = 5;
            if found then
              insert into error (f_turno, operacion, f_error, motivo)
          values (_fecha, 'Reserva', current_date, 'supera límite de reserva de turnos.');
          update error
          set nro_consultorio = subquery.nro_consultorio, 
            dni_medique = subquery.dni_medique, 
            nro_paciente = subquery.nro_paciente
        from (select t.nro_turno, a.nro_consultorio, m.dni_medique, p.nro_paciente
              from turno t, agenda a, medique m, paciente p
              where t.nro_paciente = p.nro_paciente and 
              t.dni_medique = m.dni_medique and
              t.nro_consultorio = a.nro_consultorio) as subquery
        where error.nro_error = (select (max(nro_error)) from error);
              return false;
            else
              update turno t
              set 
                nro_consultorio = a.nro_consultorio,
                nro_paciente = _nro_paciente,
                nro_obra_social_consulta = p.nro_obra_social,
                nro_afiliade_consulta = p.nro_afiliade,
                monto_paciente = c.monto_paciente,
                monto_obra_social = c.monto_obra_social,
                f_reserva = current_date,
                estado = 'Reservado'
              from cobertura c, paciente p, agenda a
              where
                t.dni_medique in (_dni_medique) and c.dni_medique in (_dni_medique)
                and t.fecha in (_fecha) and p.nro_paciente in (_nro_paciente)
                and a.dni_medique in (_dni_medique);
              return true;
            end if;
        else
          insert into error (f_turno, operacion, f_error, motivo)
          values (_fecha, 'Reserva', current_date, 'turno inexistente ó no disponible.');
          update error
          set nro_consultorio = subquery.nro_consultorio, 
            dni_medique = subquery.dni_medique, 
            nro_paciente = subquery.nro_paciente
        from (select t.nro_turno, a.nro_consultorio, m.dni_medique, p.nro_paciente
              from turno t, agenda a, medique m, paciente p
              where t.nro_paciente = p.nro_paciente and 
              t.dni_medique = m.dni_medique and
              t.nro_consultorio = a.nro_consultorio) as subquery
        where error.nro_error = (select (max(nro_error)) from error);
          return false;
        end if;
      else
        insert into error (f_turno, operacion, f_error, motivo)
        values (_fecha, 'Reserva', current_date, 'obra social de paciente no atendida por le médique.');
        update error
        set nro_consultorio = subquery.nro_consultorio, 
            dni_medique = subquery.dni_medique, 
            nro_paciente = subquery.nro_paciente
        from (select t.nro_turno, a.nro_consultorio, m.dni_medique, p.nro_paciente
              from turno t, agenda a, medique m, paciente p
              where t.nro_paciente = p.nro_paciente and 
              t.dni_medique = m.dni_medique and
              t.nro_consultorio = a.nro_consultorio) as subquery
        where error.nro_error = (select (max(nro_error)) from error);
        return false;
      end if;
    else
      insert into error (f_turno, operacion, f_error, motivo)
        values (_fecha, 'Reserva', current_date, 'nro de historia clínica no válido.');
        update error
        set nro_consultorio = subquery.nro_consultorio, 
            dni_medique = subquery.dni_medique
        from (select t.nro_turno, a.nro_consultorio, m.dni_medique
              from turno t, agenda a, medique m
              where t.dni_medique = m.dni_medique and
              t.nro_consultorio = a.nro_consultorio) as subquery
        where error.nro_error = (select (max(nro_error)) from error);
      return false;
    end if;
  else
    insert into error (f_turno, operacion, f_error, motivo)
        values (_fecha, 'Reserva', current_date, 'dni de médique no válido.');
        update error
        set nro_consultorio = subquery.nro_consultorio,
            nro_paciente = subquery.nro_paciente
        from (select t.nro_turno, a.nro_consultorio, p.nro_paciente
              from turno t, agenda a, paciente p
              where t.nro_paciente = p.nro_paciente and
              t.nro_consultorio = a.nro_consultorio) as subquery
        where error.nro_error = (select (max(nro_error)) from error);
    return false;
  end if;
end
$$ language plpgsql;
`)
	if err != nil {
		log.Fatal(err)
	}
    fmt.Printf("CREATE FUNCTION\n")
}
func reservarTurnos(db *sql.DB) {

	fmt.Printf("Ingresando consultas: solicitud_reservas\n")
	rows, err := db.Query(`select reservar_turnos(1, 44458951, '2023-11-22 08:00:00');
							select reservar_turnos(2, 44458952, '2023-11-16 11:00:00');
							select reservar_turnos(3, 44458953, '2023-11-01 07:15:00');
							select reservar_turnos(4, 44458954, '2023-12-14 06:30:00');
							select reservar_turnos(5, 44458955, '2023-12-26 12:00:00');
							select reservar_turnos(8, 44458959, '2023-12-29 16:00:00');
							select reservar_turnos(7, 44458957, '2023-12-17 15:00:00');
							select reservar_turnos(8, 44458958, '2023-11-04 17:45:00');
							select reservar_turnos(9, 44458959, '2023-12-22 14:40:00');
							select reservar_turnos(10, 44458960, '2023-11-09 14:00:00');
							select reservar_turnos(11, 44458961, '2024-01-19 18:00:00');
							select reservar_turnos(12, 44458962, '2024-11-28 19:30:00');
							select reservar_turnos(13, 44458963, '2024-12-17 23:45:00');
							select reservar_turnos(14, 44458964, '2024-12-16 20:00:00');
							select reservar_turnos(15, 44458965, '2024-01-30 21:20:00');
							select reservar_turnos(16, 44458966, '2023-11-02 02:30:00');
							select reservar_turnos(17, 44458967, '2023-12-25 04:00:00'); 
							select reservar_turnos(18, 44458968, '2023-11-25 05:00:00');
							select reservar_turnos(19, 44458969, '2023-12-05 04:30:00');
							select reservar_turnos(20, 44458970, '2023-11-08 00:00:00');
							select reservar_turnos(2, 44458952, '2023-11-02 06:00:00');
							select reservar_turnos(2, 44458955, '2023-12-12 17:00:00');
							select reservar_turnos(10, 44458952, '2023-11-02 08:00:00');
							select reservar_turnos(2, 43458952, '2023-11-09 08:00:00'); 
							select reservar_turnos(3, 44458953, '2023-11-24 12:00:00');
							select reservar_turnos(300, 44458953, '2023-11-24 11:00:00');
							select reservar_turnos(3, 44458953, '2023-11-24 10:00:00');
							select reservar_turnos(1, 44458951, '2023-11-01 08:00:00');
							select reservar_turnos(16, 44458953, '2023-11-24 11:00:00');
							select reservar_turnos(2, 44458952, '2023-11-26 08:00:00');
							select reservar_turnos(12, 44458952, '2023-11-26 06:00:00');`)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	for rows.Next(){
		var resu bool
		err = rows.Scan(&resu)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("%b", resu)
		fmt.Printf("Se insertaron todas las reservas\n")
	}
}
func anularFuncion(db *sql.DB) {
	_, err := db.Exec(`create or replace function cancelar_turnos(p_dni_medique int, p_fecha_desde date, p_fecha_hasta date, out cantidad_cancelados int)
returns integer as $$
declare
    row_data record;
begin

    update turno
    set estado = 'Cancelado'
    where dni_medique = p_dni_medique
        and fecha between p_fecha_desde and p_fecha_hasta
        and estado in ('Disponible', 'Reservado');

    get diagnostics cantidad_cancelados = row_count;
    
    insert into reprogramacion
    select t.nro_turno, p.nombre, p.apellido, p.telefono, p.email, m.nombre, m.apellido, t.estado
    from turno t, paciente p, medique m
    where t.nro_paciente = p.nro_paciente and t.dni_medique = m.dni_medique
    and estado = 'Cancelado'
    and not exists(
      select 1
      from reprogramacion r
      where r.nro_turno = t.nro_turno
    );
    
    update reprogramacion
    set estado = 'Pendiente';
    
    return;
end;
$$ language plpgsql;
`)
	if err != nil {
		log.Fatal(err)
	}
fmt.Printf("CREATE FUNCTION\n")
}
func anularTurnos(db *sql.DB) {
	fmt.Printf("se selecciono el medico 44458952 y los turnos a cancelar son del 2023-11-25 al 2023-11-27\n")
	rows, err := db.Query("select cancelar_turnos(44458952, '2023-11-25', '2023-11-27');")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	for rows.Next(){
		var resu int
		err = rows.Scan(&resu)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("se cancelaron %d turnos\n", resu)
	}
}
func atencionFuncion(db *sql.DB) {
	_, err := db.Exec(`create or replace function atender_turno(nro_turno_hg int) returns boolean as $$
declare
    v_estado_turno char(10);
    v_fecha_turno date;
begin
 
 select estado, fecha into v_estado_turno, v_fecha_turno
    from turno
    where  nro_turno = nro_turno_hg;
    if not found then
        insert into error (f_turno, operacion, f_error, motivo)
        values (current_date, 'Atención', current_date, 'nro de turno no válido');
        
        return false;
    end if;
    if v_estado_turno != 'Reservado' then
        insert into error (f_turno, operacion, f_error, motivo)
        values (current_date, 'Atención', current_date, 'turno no reservado');
        return false;
    end if;
    if v_fecha_turno != current_date then
        insert into error (f_turno, operacion, f_error, motivo)
        values (current_date, 'Atención', current_date, 'turno no corresponde a la fecha del día');
        update error
        set nro_consultorio = subquery.nro_consultorio, 
            dni_medique = subquery.dni_medique, 
            nro_paciente = subquery.nro_paciente
        from (select t.nro_turno, a.nro_consultorio, m.dni_medique, p.nro_paciente
              from turno t, agenda a, medique m, paciente p
              where t.nro_paciente = p.nro_paciente and 
              t.dni_medique = m.dni_medique and
              t.nro_consultorio = a.nro_consultorio) as subquery
        where error.nro_error = (select (max(nro_error)) from error);     
        return false;
    end if;
    update turno
    set estado = 'Atendido'
    where nro_turno = nro_turno_hg;
    return true;
end
$$ language plpgsql;
`)
	if err != nil {
		log.Fatal(err)
	}
fmt.Printf("CREATE FUNCTION\n")
}

func atencionTurnos(db *sql.DB) {
	fmt.Printf("Se atendieron los turnos 292 y 296 ")
	rows, err := db.Query(`select atender_turno(292);
							select atender_turno(296);`)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	for rows.Next(){
		var resu bool
		err = rows.Scan(&resu)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("%b", resu)
	}
}

func liqObraSocialFuncion(db *sql.DB) {
	_, err := db.Exec(`create or replace function generar_liquidacion_mensual(p_nro_obra_social INT, p_desde DATE, p_hasta DATE)
returns void as $$
declare
begin
    
    insert into liquidacion_cabecera (nro_obra_social, desde, hasta, total)
    values(p_nro_obra_social, p_desde, p_hasta, 0);
    
    update liquidacion_cabecera
	set total = (select sum(t.monto_obra_social) from turno t
	where estado = 'Atendido' 
	and t.nro_obra_social_consulta = liquidacion_cabecera.nro_obra_social
	group by extract(month from t.fecha));

   update turno
    set estado = 'Liquidado'
    where fecha between p_desde and p_hasta
    and nro_obra_social_consulta = p_nro_obra_social
    and estado = 'Atendido';

    insert into liquidacion_detalle (nro_liquidacion)
    select l.nro_liquidacion
    from turno t, liquidacion_cabecera l
    where estado = 'Liquidado'
    and not exists(
      select 1
      from liquidacion_detalle d
      where d.nro_liquidacion = l.nro_liquidacion
    );
    
    update liquidacion_detalle
    set f_atencion = subquery.fecha,
        nro_afiliade = subquery.nro_afiliade,
        dni_paciente = subquery.dni_paciente,
        nombre_paciente = subquery.nombre_paciente,
        apellido_paciente = subquery.apellido_paciente,
        dni_medique = subquery.dni_medique,
        nombre_medique = subquery.nombre,
        apellido_medique = subquery.apellido,
        especialidad = subquery.especialidad,
        monto = subquery.monto_obra_social
    from (select distinct t.fecha, p.nro_afiliade, p.dni_paciente,
          p.nombre as nombre_paciente, p.apellido as apellido_paciente, m.dni_medique, m.nombre, m.apellido,
          m.especialidad, t.monto_obra_social
          from turno t, paciente p, medique m
          where t.nro_afiliade_consulta = p.nro_afiliade and
                t.dni_medique = m.dni_medique) as subquery, turno
    where turno.fecha between p_desde and p_hasta;
end
$$ language plpgsql;
`)
if err != nil {
	log.Fatal(err)
}
fmt.Printf("CREATE FUNCTION\n")

}
func generarLiqObraSocial(db *sql.DB) {
	fmt.Printf("Se genero la liquidacion del mes de noviembre ")
	_ , err := db.Query(`select generar_liquidacion_mensual(721, '2023-11-01', '2023-12-01');
							select generar_liquidacion_mensual(722, '2023-11-01', '2023-12-01');
							select generar_liquidacion_mensual(723, '2023-11-01', '2023-12-01');`)
	if err != nil {
		log.Fatal(err)
	}
}

func envioEmail(db *sql.DB) {
	_, err := db.Exec(`create or replace function enviar_email()
returns trigger as $$
declare
  dia cursor for select extract(day from current_date);
begin
  case new.estado 
  when 'Reservado' then
    insert into envio_email(f_generacion, asunto, cuerpo, f_envio, estado) 
    values (current_date, 'Reserva de turno', 'El turno fue reservado correctamente', null, 'Pendiente');
    if (date_trunc('day', current_date) != date_trunc('day', current_date - interval '1 day')) then
      if (date_trunc('day', new.fecha) = date_trunc('day', current_timestamp) + interval '2 days') then
        insert into envio_email(f_generacion, asunto, f_envio, estado)
        values (current_date, 'Recordatorio de turno', null, 'Pendiente');
        update envio_email set cuerpo = _cuerpo
        from (select t.nro_turno, t.fecha, t.dni_medique from turno t 
        where (date_trunc('day', t.fecha) = date_trunc('day', current_timestamp) + interval '2 days')) as _cuerpo
        where asunto = 'Recordatorio de turno';
      end if;
    end if;
    if exists(select t.fecha from turno t where (date_trunc('day', fecha) = date_trunc('day', current_timestamp - interval '1 day'))) then
      insert into envio_email(f_generacion, asunto, f_envio, estado)
      values (current_date, 'Perdida de turno reservado', null, 'Pendiente');
      update envio_email set cuerpo = cuerpo
        from (select t.nro_turno, t.fecha, t.dni_medique from turno t 
        where (date_trunc('day', t.fecha) = date_trunc('day', current_timestamp) - interval '1 day')) as cuerpo
        where asunto = 'Perdida de turno reservado';
    end if;
  when 'Cancelado' then
    insert into envio_email(f_generacion, asunto, cuerpo, f_envio, estado) 
    values (current_date, 'Cancelacion de turno', 'El turno fue CANCELADO', null, 'Pendiente');
  else 
    return new;
  end case;
  return new;
end
$$ language plpgsql;

create trigger email_reserva_turno
after update of estado on turno
for each row
when (new.estado = 'Reservado')
execute function enviar_email();

create trigger email_cancelar_turno
after update of estado on turno
for each row
when (new.estado = 'Cancelado')
execute function enviar_email();

create trigger email_recordatorio_turno
after update on turno
for each row
when (date_trunc('day', new.fecha) = date_trunc('day', current_timestamp) + interval '2 days')
execute function enviar_email();

create trigger email_perdida_turno
after update on turno
for each row
when (date_trunc('day', current_date) != date_trunc('day', current_date - interval '1 day'))
execute function enviar_email();
`)
	if err != nil {
		log.Fatal(err)
	}
fmt.Printf("CREATE FUNCTION\n")
}
func CreateUpdate(db *bolt.DB, bucketName string, key []byte, val []byte) error {
    tx, err := db.Begin(true)
    if err != nil {
        return err
    }
    defer tx.Rollback()

    b, _ := tx.CreateBucketIfNotExists([]byte(bucketName))

    err = b.Put(key, val)
    if err != nil {
        return err
    }

    if err := tx.Commit(); err != nil {
        return err
    }

    return nil
}
func ReadUnique(db *bolt.DB, bucketName string, key []byte) ([]byte, error) {
    var buf []byte

    err := db.View(func(tx *bolt.Tx) error {
        b := tx.Bucket([]byte(bucketName))
        buf = b.Get(key)
        return nil
    })

    return buf, err
}
func boltStart(db *sql.DB){
    var opcion int
	db1, err := bolt.Open("mi_bolt.db", 0600, nil)
    if err != nil {
        log.Fatal(err)
    }
    defer db1.Close()
	for {
		fmt.Printf("\nMenu > Postgres to BoldDB \n-ObraSocial : 1\n-Consultorios : 2\n-Pacientes : 3\n-Mediques : 4\n-Turno : 5\n-Volver : 0\n")
		fmt.Scan(&opcion)
		if opcion == 1 {
			boltObraSocial(db, db1)
		}
		if opcion == 2 {
			boltConsultorio(db, db1)		
		}
		if opcion == 3 {
			boltPaciente(db, db1)	
		}
		if opcion == 4 {
			boltMedique(db, db1)
		}
		if opcion == 5{
			boltTurno(db, db1)
		}
		if opcion == 0 {
			fmt.Printf("Volviendo..")
			break
		}
	}
}

func boltObraSocial(db *sql.DB, db1 *bolt.DB ){
	rows, err := db.Query("select * from obra_social")
    if err != nil {
        log.Fatal(err)
    }
	defer rows.Close()
	for rows.Next() {
		var nro_obra_social int
		var nombre string
		var contacto_nombre string
		var contacto_apellido string
		var contacto_telefono string
		var contacto_email string
		err = rows.Scan(&nro_obra_social, &nombre, &contacto_nombre, &contacto_apellido, &contacto_telefono, &contacto_email)
		if err != nil {
			panic(err)
		}
		x := Obra_social{nro_obra_social, nombre, contacto_nombre, contacto_apellido, contacto_telefono, contacto_email}
		data, err := json.Marshal(x)
		if err != nil {
			log.Fatal(err)
		}
		CreateUpdate(db1, "obra_social", []byte(strconv.Itoa(x.Nro_obra_social)), data)
		resultado, err := ReadUnique(db1, "obra_social", []byte(strconv.Itoa(x.Nro_obra_social)))
		fmt.Printf("%s\n", resultado)
	}
	fmt.Printf("--> true\n")
}
func boltConsultorio(db *sql.DB, db1 *bolt.DB ){
	rows, err := db.Query("select * from consultorio")
    if err != nil {
        log.Fatal(err)
    }
	defer rows.Close()
	for rows.Next() {
		var nro_consultorio int
		var nombre string
		var domicilio string
		var codigo_postal string
		var telefono string
		err = rows.Scan(&nro_consultorio, &nombre, &domicilio, &codigo_postal, &telefono)
		if err != nil {
			panic(err)
		}
		x := Consultorio{nro_consultorio, nombre, domicilio, codigo_postal, telefono}
		data, err := json.Marshal(x)
		if err != nil {
			log.Fatal(err)
		}
		CreateUpdate(db1, "consultorio", []byte(strconv.Itoa(x.Nro_consultorio)), data)
		resultado, err := ReadUnique(db1, "consultorio", []byte(strconv.Itoa(x.Nro_consultorio)))
		fmt.Printf("%s\n", resultado)
	}
	fmt.Printf("--> true\n")
}
func boltPaciente(db *sql.DB, db1 *bolt.DB ){
	rows, err := db.Query("select * from paciente")
    if err != nil {
        log.Fatal(err)
    }
	defer rows.Close()
	for rows.Next() {
		var nro_paciente int
		var nombre string
		var apellido string
		var dni_paciente int
		var f_nac string
		var nro_obra_social int
		var nro_afiliade int
		var domicilio string
		var telefono string
		var email string
		err = rows.Scan(&nro_paciente, &nombre, &apellido, &dni_paciente, &f_nac, &nro_obra_social, &nro_afiliade, &domicilio, &telefono, &email)
		if err != nil {
			panic(err)
		}
		x := Paciente{nro_paciente, nombre, apellido, dni_paciente, f_nac, nro_obra_social, nro_afiliade, domicilio, telefono, email}
		data, err := json.Marshal(x)
		if err != nil {
			panic(err)
		}
		CreateUpdate(db1, "paciente", []byte(strconv.Itoa(x.Nro_paciente)), data)
		resultado, err := ReadUnique(db1, "paciente", []byte(strconv.Itoa(x.Nro_paciente)))
		fmt.Printf("%s\n", resultado)
	}
	fmt.Printf("--> true\n")
}
func boltMedique(db *sql.DB, db1 *bolt.DB ){
	rows, err := db.Query("select * from medique")
    if err != nil {
        log.Fatal(err)
    }
	defer rows.Close()
	for rows.Next() {
		var dni_medique int
		var nombre string
		var apellido string
		var especialidad string
		var monto_consulta_privada float64
		var telefono string
		err = rows.Scan(&dni_medique, &nombre, &apellido, &especialidad, &monto_consulta_privada, &telefono)
		if err != nil {
			panic(err)
		}
		x := Medique{dni_medique, nombre, apellido, especialidad, monto_consulta_privada, telefono}
		data, err := json.Marshal(x)
		if err != nil {
			log.Fatal(err)
		}
		CreateUpdate(db1, "medique", []byte(strconv.Itoa(x.Dni_medique)), data)
		resultado, err := ReadUnique(db1, "medique", []byte(strconv.Itoa(x.Dni_medique)))
		fmt.Printf("%s\n", resultado)
	}
	fmt.Printf("--> true\n")
}
func boltTurno(db *sql.DB, db1 *bolt.DB ){
	rows, err := db.Query(`select nro_turno, fecha, nro_consultorio, dni_medique, nro_paciente, nro_obra_social_consulta, nro_afiliade_consulta, monto_paciente, monto_obra_social, f_reserva, estado 
							from (select *, row_number() over (partition by dni_medique order by nro_turno desc) as row_num from turno t where t.estado='Reservado') subquery where row_num <= 3;`)
    if err != nil {
        log.Fatal(err)
    }
	defer rows.Close()
	for rows.Next() {
		var nro_turno int
		var fecha string
		var nro_consultorio int
		var dni_medique int
		var nro_paciente int
		var nro_obra_social_consulta int
		var nro_afiliade_consulta int
		var monto_paciente float64
		var monto_obra_social float64
		var f_reserva string
		var estado string
		err = rows.Scan(&nro_turno, &fecha, &nro_consultorio, &dni_medique, &nro_paciente, &nro_obra_social_consulta, &nro_afiliade_consulta, &monto_paciente, &monto_obra_social, &f_reserva, &estado)
		if err != nil {
			panic(err)
		}
		x := Turno{nro_turno, fecha, nro_consultorio, dni_medique, nro_paciente, nro_obra_social_consulta, nro_afiliade_consulta, monto_paciente, monto_obra_social, f_reserva, estado}
		data, err := json.Marshal(x)
		if err != nil {
			panic(err)
		}
		CreateUpdate(db1, "turno", []byte(strconv.Itoa(x.Nro_turno)), data)
		resultado, err := ReadUnique(db1, "turno", []byte(strconv.Itoa(x.Nro_turno)))
		fmt.Printf("%s\n", resultado)
	}
	fmt.Printf("--> true\n")
}
type Obra_social struct {
	Nro_obra_social int
	Nombre string
	Contacto_nombre string
	Contacto_apellido string
	Contacto_telefono string
	Contacto_email string
}

type Consultorio struct {
	Nro_consultorio int
	Nombre string
	Domicilio string
	Codigo_postal string
	Telefono string
}
type Paciente struct {
	Nro_paciente int
	Nombre string
	Apellido string
	Dni_paciente int
	F_nac string
	Nro_obra_social int
	Nro_afiliade int
	Domicilio string
	Telefono string
	Email string
}
type Medique struct {
	Dni_medique int
	Nombre string
	Apellido string
	Especialidad string
	Monto_consulta_privada float64
	Telefono string
}
type Turno struct {
	Nro_turno int
	Fecha string
	Nro_consultorio int //`no supe que querias que renombrara con esto`
	Dni_medique int
	Nro_paciente int
	Nro_obra_social_consulta int
	Nro_afiliade_consulta int
	Monto_paciente float64
	Monto_obra_social float64
	F_reserva string
	Estado string
}
