CREATE DATABASE natsservice;
\c natsservice;

CREATE TABLE delivery (
                          id SERIAL PRIMARY KEY ,
                          name text not null ,
                          phone text not null ,
                          zip text not null ,
                          city text not null ,
                          address text not null ,
                          region text not null ,
                          email text not null
);

CREATE TABLE payment (
                         id SERIAL PRIMARY KEY ,
                         transaction text not null ,
                         request_id text null ,
                         currency text not null ,
                         provider text not null ,
                         amount int not null ,
                         payment_dt int not null ,
                         bank text not null ,
                         delivery_cost int not null ,
                         goods_total int not null ,
                         custom_fee int null
);

CREATE TABLE items (
                       id SERIAL PRIMARY KEY ,
                       chrt_id int unique not null ,
                       track_number text not null ,
                       price int not null ,
                       rid text not null ,
                       name text not null ,
                       sale int not null ,
                       size text not null ,
                       total_price int not null ,
                       nm_id int not null ,
                       brand text not null ,
                       status int not null
);

CREATE TABLE wb_order (
                          id SERIAL PRIMARY KEY,

                          order_uid text unique ,
                          track_number text not null,
                          entry text not null,
                          delivery_id INT not null,
                          payment_id INT not null,
                          locale text not null,
                          internal_signature TEXT,
                          customer_id text not null unique ,
                          delivery_service text not null,
                          shardkey text not null ,
                          sm_id int not null ,
                          date_created date not null ,
                          oof_shard text not null,

                          CONSTRAINT fk_delivery
                              FOREIGN KEY (delivery_id)
                                  REFERENCES delivery(id)
                                  ON DELETE CASCADE ,

                          CONSTRAINT fk_payment
                              FOREIGN KEY (payment_id)
                                  REFERENCES payment(id)
                                  ON DELETE CASCADE
);

CREATE TABLE order_items (
                             id SERIAL PRIMARY KEY ,
                             order_uid text not null ,
                             chrt_id int not null,
                             CONSTRAINT fk_order_uid
                                 FOREIGN KEY (order_uid)
                                     REFERENCES wb_order(order_uid)
                                     ON DELETE CASCADE,
                             CONSTRAINT fk_chrt_id
                                 FOREIGN KEY (chrt_id)
                                     REFERENCES items(chrt_id)
                                     ON DELETE CASCADE
);
