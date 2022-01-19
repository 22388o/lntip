CREATE TRIGGER update_balance_trigger
AFTER INSERT ON tips
FOR EACH ROW
BEGIN
    UPDATE users SET balance = balance + NEW.amount WHERE users.id = NEW.to_user_id;
    UPDATE users SET balance = balance - NEW.amount WHERE users.id = NEW.user_id;
END;
