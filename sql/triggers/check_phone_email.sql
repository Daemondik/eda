DO
$$
    BEGIN
        IF EXISTS (SELECT 1 FROM pg_trigger WHERE tgname = 'check_phone_email') THEN
            DROP TRIGGER check_phone_email ON users;
        END IF;

        CREATE OR REPLACE FUNCTION check_phone_email()
            RETURNS TRIGGER AS
        $func$
        DECLARE
            valid_email_pattern TEXT := '^[A-Za-z0-9._%-]+@[A-Za-z0-9.-]+\.[A-Za-z]{2,4}$';
            valid_phone_pattern TEXT := '^7\d{10}$';
        BEGIN
            IF NEW.phone IS NOT NULL AND NEW.email IS NULL THEN
                IF EXISTS (SELECT 1 FROM users WHERE phone = NEW.phone) THEN
                    RAISE EXCEPTION 'Phone number already exists';
                END IF;
                IF NOT NEW.phone ~ valid_phone_pattern THEN
                    RAISE EXCEPTION 'Invalid phone number';
                END IF;
            ELSIF NEW.email IS NOT NULL AND NEW.phone IS NULL THEN
                IF EXISTS (SELECT 1 FROM users WHERE email = NEW.email) THEN
                    RAISE EXCEPTION 'Email already exists';
                END IF;
                IF NOT NEW.email ~ valid_email_pattern THEN
                    RAISE EXCEPTION 'Invalid email';
                END IF;
            END IF;
            RETURN NEW;
        END;
        $func$ LANGUAGE plpgsql;

        CREATE TRIGGER check_phone_email
            BEFORE INSERT OR UPDATE
            ON users
            FOR EACH ROW
        EXECUTE PROCEDURE check_phone_email();
    END;
$$;