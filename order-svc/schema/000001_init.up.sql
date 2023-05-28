-- SET TIME ZONE 'Europe/Moscow';
SET TIME ZONE 0;

CREATE TABLE dishes
(
    id           UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name         VARCHAR(100)   NOT NULL,
    description  TEXT,
    price        DECIMAL(10, 2) NOT NULL,
--     price        MONEY        NOT NULL,
    quantity     INT            NOT NULL,
    --     dish can't be deleted, it can only be marked as unavailable
    is_available BOOLEAN        NOT NULL,
    created_at   TIMESTAMP        DEFAULT CURRENT_TIMESTAMP,
    updated_at   TIMESTAMP        DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT name_unique UNIQUE ("name")
);

CREATE OR REPLACE FUNCTION update_changetimestamp_column()
    RETURNS TRIGGER AS
$$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ language 'plpgsql';

CREATE TRIGGER update_dish_changetimestamp
    BEFORE UPDATE
    ON dishes
    FOR EACH ROW
EXECUTE PROCEDURE
    update_changetimestamp_column();

CREATE TABLE orders
(
    id               UUID PRIMARY KEY     DEFAULT gen_random_uuid(),
    user_id          UUID        NOT NULL,
    status           VARCHAR(50) NOT NULL DEFAULT 'cancelled',
    special_requests TEXT,
    ready_at         TIMESTAMP   NOT NULL,
    created_at       TIMESTAMP            DEFAULT CURRENT_TIMESTAMP,
    updated_at       TIMESTAMP            DEFAULT CURRENT_TIMESTAMP
);

CREATE TRIGGER update_orders_changetimestamp
    BEFORE UPDATE
    ON orders
    FOR EACH ROW
EXECUTE PROCEDURE
    update_changetimestamp_column();

-- many to one
-- changedPRIMARY KEY to (order_id, dish_id) from
-- id       UUID PRIMARY KEY,
CREATE TABLE order_dish
(
    order_id UUID           NOT NULL,
    dish_id  UUID           NOT NULL,
    quantity INT            NOT NULL,
    price    DECIMAL(10, 2) NOT NULL,
    FOREIGN KEY (order_id) REFERENCES orders (id),
    FOREIGN KEY (dish_id) REFERENCES dishes (id),
    CONSTRAINT order_dish_u PRIMARY KEY (order_id, dish_id)
);

INSERT INTO dishes (id, "name", description, price, quantity, is_available)
VALUES ('90c193c8-9f23-4c3d-bcff-6963bbc37034' ,'Carbonara', 'Traditional pasta (Spaghetti or Penne) with eggs, pecorino romano cheese, guanciale and black pepper', 15.99,
        10, true),
       ('828b72c6-a397-445b-af92-62cebe548491', 'Caprese Salad', 'Tomatoes, mozzarella, and fresh basil salad', 9.99, 3, true),
       ('3f5382d4-4a4f-4750-841e-8bf120195891', 'Tiramisu', 'Coffee-flavored dessert made with ladyfingers, mascarpone cheese, fresh eggs and cocoa powder',
        6.99, 10, true);

INSERT INTO dishes ("name", description, price, quantity, is_available)
VALUES
        ('Margherita Pizza', 'San marzano tomato sauce, mozzarella and basil', 12.99, 10, true),
       ('Lasagna', 'Oven-baked pasta dish with layers of pasta, ragu alla bolognese, and cheese', 16.99, 1, true),
       ('Spaghetti al ragù (alla Bolognese)', 'Spaghetti pasta with a meaty tomato bolognese ragu', 14.99, 5, true),
       ('Bruschetta', 'Toasted bread topped with garlic, tomatoes, and basil', 7.99, 20, true),
       ('Risotto', 'Creamy rice dish cooked with broth, porcini mushrooms and pancetta', 18.99, 10, true),
       ('Minestrone Soup', 'Vegetable soup with pasta or rice', 8.99, 30, true);

