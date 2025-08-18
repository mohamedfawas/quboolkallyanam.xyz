DO
$$
BEGIN
	IF NOT EXISTS (SELECT 1 FROM pg_database WHERE datname = 'qubool_kallyanam_auth') THEN
		CREATE DATABASE qubool_kallyanam_auth;
	END IF;

	IF NOT EXISTS (SELECT 1 FROM pg_database WHERE datname = 'qubool_kallyanam_user') THEN
		CREATE DATABASE qubool_kallyanam_user;
	END IF;

	IF NOT EXISTS (SELECT 1 FROM pg_database WHERE datname = 'qubool_kallyanam_chat') THEN
		CREATE DATABASE qubool_kallyanam_chat;
	END IF;

	IF NOT EXISTS (SELECT 1 FROM pg_database WHERE datname = 'qubool_kallyanam_payment') THEN
		CREATE DATABASE qubool_kallyanam_payment;
	END IF;
END
$$;