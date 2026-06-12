--
-- PostgreSQL database dump
--

\restrict WziVqJckDY8ZIhPD6Rfejs92g9wO5hFJTAfatPag4ldGQEsTfroGlL3IyUWmeUQ

-- Dumped from database version 15.18
-- Dumped by pg_dump version 16.14 (Ubuntu 16.14-0ubuntu0.24.04.1)

SET statement_timeout = 0;
SET lock_timeout = 0;
SET idle_in_transaction_session_timeout = 0;
SET client_encoding = 'UTF8';
SET standard_conforming_strings = on;
SELECT pg_catalog.set_config('search_path', '', false);
SET check_function_bodies = false;
SET xmloption = content;
SET client_min_messages = warning;
SET row_security = off;

SET default_tablespace = '';

SET default_table_access_method = heap;

--
-- Name: migrations; Type: TABLE; Schema: public; Owner: bagogo
--

CREATE TABLE public.migrations (
    id character varying(255) NOT NULL
);


ALTER TABLE public.migrations OWNER TO bagogo;

--
-- Name: order_sequences; Type: TABLE; Schema: public; Owner: bagogo
--

CREATE TABLE public.order_sequences (
    id uuid DEFAULT gen_random_uuid() NOT NULL,
    date date NOT NULL,
    last_sequence bigint NOT NULL,
    created_at timestamp with time zone,
    updated_at timestamp with time zone,
    deleted_at timestamp with time zone
);


ALTER TABLE public.order_sequences OWNER TO bagogo;

--
-- Name: orders; Type: TABLE; Schema: public; Owner: bagogo
--

CREATE TABLE public.orders (
    id uuid DEFAULT gen_random_uuid() NOT NULL,
    created_at timestamp with time zone,
    created_by uuid,
    updated_at timestamp with time zone,
    updated_by uuid,
    deleted_at timestamp with time zone,
    user_id uuid,
    grand_total bigint NOT NULL,
    status character varying(20) DEFAULT 'draft'::character varying,
    order_number character varying(50)
);


ALTER TABLE public.orders OWNER TO bagogo;

--
-- Name: products; Type: TABLE; Schema: public; Owner: bagogo
--

CREATE TABLE public.products (
    id uuid DEFAULT gen_random_uuid() NOT NULL,
    created_at timestamp with time zone,
    created_by uuid,
    updated_at timestamp with time zone,
    updated_by uuid,
    deleted_at timestamp with time zone,
    name text NOT NULL,
    description text,
    slug text,
    price bigint NOT NULL,
    stock bigint DEFAULT 0,
    unit character varying(10),
    is_active boolean DEFAULT true,
    store_id uuid
);


ALTER TABLE public.products OWNER TO bagogo;

--
-- Name: stores; Type: TABLE; Schema: public; Owner: bagogo
--

CREATE TABLE public.stores (
    id uuid DEFAULT gen_random_uuid() NOT NULL,
    created_at timestamp with time zone,
    created_by uuid,
    updated_at timestamp with time zone,
    updated_by uuid,
    deleted_at timestamp with time zone,
    name text NOT NULL,
    description text,
    is_active boolean DEFAULT true,
    user_id uuid
);


ALTER TABLE public.stores OWNER TO bagogo;

--
-- Name: users; Type: TABLE; Schema: public; Owner: bagogo
--

CREATE TABLE public.users (
    id uuid DEFAULT gen_random_uuid() NOT NULL,
    created_at timestamp with time zone,
    created_by uuid,
    updated_at timestamp with time zone,
    updated_by uuid,
    deleted_at timestamp with time zone,
    username text NOT NULL,
    email text NOT NULL,
    password text NOT NULL,
    role character varying(20) DEFAULT 'buyer'::character varying,
    is_active boolean DEFAULT true,
    picture character varying(500)
);


ALTER TABLE public.users OWNER TO bagogo;

--
-- Name: migrations migrations_pkey; Type: CONSTRAINT; Schema: public; Owner: bagogo
--

ALTER TABLE ONLY public.migrations
    ADD CONSTRAINT migrations_pkey PRIMARY KEY (id);


--
-- Name: order_sequences order_sequences_pkey; Type: CONSTRAINT; Schema: public; Owner: bagogo
--

ALTER TABLE ONLY public.order_sequences
    ADD CONSTRAINT order_sequences_pkey PRIMARY KEY (id);


--
-- Name: orders orders_pkey; Type: CONSTRAINT; Schema: public; Owner: bagogo
--

ALTER TABLE ONLY public.orders
    ADD CONSTRAINT orders_pkey PRIMARY KEY (id);


--
-- Name: products products_pkey; Type: CONSTRAINT; Schema: public; Owner: bagogo
--

ALTER TABLE ONLY public.products
    ADD CONSTRAINT products_pkey PRIMARY KEY (id);


--
-- Name: stores stores_pkey; Type: CONSTRAINT; Schema: public; Owner: bagogo
--

ALTER TABLE ONLY public.stores
    ADD CONSTRAINT stores_pkey PRIMARY KEY (id);


--
-- Name: stores uni_stores_name; Type: CONSTRAINT; Schema: public; Owner: bagogo
--

ALTER TABLE ONLY public.stores
    ADD CONSTRAINT uni_stores_name UNIQUE (name);


--
-- Name: users uni_users_email; Type: CONSTRAINT; Schema: public; Owner: bagogo
--

ALTER TABLE ONLY public.users
    ADD CONSTRAINT uni_users_email UNIQUE (email);


--
-- Name: users uni_users_username; Type: CONSTRAINT; Schema: public; Owner: bagogo
--

ALTER TABLE ONLY public.users
    ADD CONSTRAINT uni_users_username UNIQUE (username);


--
-- Name: users users_pkey; Type: CONSTRAINT; Schema: public; Owner: bagogo
--

ALTER TABLE ONLY public.users
    ADD CONSTRAINT users_pkey PRIMARY KEY (id);


--
-- Name: idx_order_sequences_date; Type: INDEX; Schema: public; Owner: bagogo
--

CREATE UNIQUE INDEX idx_order_sequences_date ON public.order_sequences USING btree (date);


--
-- Name: idx_order_sequences_deleted_at; Type: INDEX; Schema: public; Owner: bagogo
--

CREATE INDEX idx_order_sequences_deleted_at ON public.order_sequences USING btree (deleted_at);


--
-- Name: idx_orders_deleted_at; Type: INDEX; Schema: public; Owner: bagogo
--

CREATE INDEX idx_orders_deleted_at ON public.orders USING btree (deleted_at);


--
-- Name: idx_orders_grand_total; Type: INDEX; Schema: public; Owner: bagogo
--

CREATE INDEX idx_orders_grand_total ON public.orders USING btree (grand_total);


--
-- Name: idx_orders_order_number; Type: INDEX; Schema: public; Owner: bagogo
--

CREATE UNIQUE INDEX idx_orders_order_number ON public.orders USING btree (order_number);


--
-- Name: idx_orders_status; Type: INDEX; Schema: public; Owner: bagogo
--

CREATE INDEX idx_orders_status ON public.orders USING btree (status);


--
-- Name: idx_orders_user_id; Type: INDEX; Schema: public; Owner: bagogo
--

CREATE INDEX idx_orders_user_id ON public.orders USING btree (user_id);


--
-- Name: idx_products_deleted_at; Type: INDEX; Schema: public; Owner: bagogo
--

CREATE INDEX idx_products_deleted_at ON public.products USING btree (deleted_at);


--
-- Name: idx_products_is_active; Type: INDEX; Schema: public; Owner: bagogo
--

CREATE INDEX idx_products_is_active ON public.products USING btree (is_active);


--
-- Name: idx_products_price; Type: INDEX; Schema: public; Owner: bagogo
--

CREATE INDEX idx_products_price ON public.products USING btree (price);


--
-- Name: idx_products_slug; Type: INDEX; Schema: public; Owner: bagogo
--

CREATE UNIQUE INDEX idx_products_slug ON public.products USING btree (slug);


--
-- Name: idx_products_stock; Type: INDEX; Schema: public; Owner: bagogo
--

CREATE INDEX idx_products_stock ON public.products USING btree (stock);


--
-- Name: idx_products_store_id; Type: INDEX; Schema: public; Owner: bagogo
--

CREATE INDEX idx_products_store_id ON public.products USING btree (store_id);


--
-- Name: idx_stores_deleted_at; Type: INDEX; Schema: public; Owner: bagogo
--

CREATE INDEX idx_stores_deleted_at ON public.stores USING btree (deleted_at);


--
-- Name: idx_stores_is_active; Type: INDEX; Schema: public; Owner: bagogo
--

CREATE INDEX idx_stores_is_active ON public.stores USING btree (is_active);


--
-- Name: idx_stores_user_id; Type: INDEX; Schema: public; Owner: bagogo
--

CREATE INDEX idx_stores_user_id ON public.stores USING btree (user_id);


--
-- Name: idx_users_deleted_at; Type: INDEX; Schema: public; Owner: bagogo
--

CREATE INDEX idx_users_deleted_at ON public.users USING btree (deleted_at);


--
-- Name: orders fk_orders_user; Type: FK CONSTRAINT; Schema: public; Owner: bagogo
--

ALTER TABLE ONLY public.orders
    ADD CONSTRAINT fk_orders_user FOREIGN KEY (user_id) REFERENCES public.users(id);


--
-- Name: products fk_products_store; Type: FK CONSTRAINT; Schema: public; Owner: bagogo
--

ALTER TABLE ONLY public.products
    ADD CONSTRAINT fk_products_store FOREIGN KEY (store_id) REFERENCES public.stores(id);


--
-- Name: stores fk_stores_user; Type: FK CONSTRAINT; Schema: public; Owner: bagogo
--

ALTER TABLE ONLY public.stores
    ADD CONSTRAINT fk_stores_user FOREIGN KEY (user_id) REFERENCES public.users(id);


--
-- PostgreSQL database dump complete
--

\unrestrict WziVqJckDY8ZIhPD6Rfejs92g9wO5hFJTAfatPag4ldGQEsTfroGlL3IyUWmeUQ

