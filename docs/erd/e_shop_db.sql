--
-- PostgreSQL database dump
--

\restrict aJmBIO3L7VKA7TlJsdaeP6fiTN3gHuhh0Wi9Cy9jzvZPgxMtEHcCdGzpTOfSi6c

-- Dumped from database version 15.17
-- Dumped by pg_dump version 16.13 (Ubuntu 16.13-0ubuntu0.24.04.1)

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
-- Name: order_items; Type: TABLE; Schema: public; Owner: admin
--

CREATE TABLE public.order_items (
    id uuid DEFAULT gen_random_uuid() NOT NULL,
    created_at timestamp with time zone,
    created_by uuid,
    updated_at timestamp with time zone,
    updated_by uuid,
    deleted_at timestamp with time zone,
    order_id uuid,
    store_id uuid,
    product_id uuid,
    quantity bigint NOT NULL,
    price bigint NOT NULL,
    sub_total bigint NOT NULL
);


ALTER TABLE public.order_items OWNER TO admin;

--
-- Name: orders; Type: TABLE; Schema: public; Owner: admin
--

CREATE TABLE public.orders (
    id uuid DEFAULT gen_random_uuid() NOT NULL,
    created_at timestamp with time zone,
    created_by uuid,
    updated_at timestamp with time zone,
    updated_by uuid,
    deleted_at timestamp with time zone,
    user_id uuid,
    status character varying(20) DEFAULT 'draft'::character varying,
    grand_total bigint NOT NULL
);


ALTER TABLE public.orders OWNER TO admin;

--
-- Name: products; Type: TABLE; Schema: public; Owner: admin
--

CREATE TABLE public.products (
    id uuid DEFAULT gen_random_uuid() NOT NULL,
    created_at timestamp with time zone,
    created_by uuid,
    updated_at timestamp with time zone,
    updated_by uuid,
    deleted_at timestamp with time zone,
    name text NOT NULL,
    slug text,
    price bigint NOT NULL,
    stock bigint DEFAULT 0,
    store_id uuid,
    description text,
    unit character varying(10),
    is_active boolean DEFAULT true
);


ALTER TABLE public.products OWNER TO admin;

--
-- Name: stores; Type: TABLE; Schema: public; Owner: admin
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
    user_id uuid,
    is_active boolean DEFAULT true
);


ALTER TABLE public.stores OWNER TO admin;

--
-- Name: users; Type: TABLE; Schema: public; Owner: admin
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
    is_active boolean DEFAULT true
);


ALTER TABLE public.users OWNER TO admin;

--
-- Name: order_items order_items_pkey; Type: CONSTRAINT; Schema: public; Owner: admin
--

ALTER TABLE ONLY public.order_items
    ADD CONSTRAINT order_items_pkey PRIMARY KEY (id);


--
-- Name: orders orders_pkey; Type: CONSTRAINT; Schema: public; Owner: admin
--

ALTER TABLE ONLY public.orders
    ADD CONSTRAINT orders_pkey PRIMARY KEY (id);


--
-- Name: products products_pkey; Type: CONSTRAINT; Schema: public; Owner: admin
--

ALTER TABLE ONLY public.products
    ADD CONSTRAINT products_pkey PRIMARY KEY (id);


--
-- Name: stores stores_pkey; Type: CONSTRAINT; Schema: public; Owner: admin
--

ALTER TABLE ONLY public.stores
    ADD CONSTRAINT stores_pkey PRIMARY KEY (id);


--
-- Name: stores uni_stores_name; Type: CONSTRAINT; Schema: public; Owner: admin
--

ALTER TABLE ONLY public.stores
    ADD CONSTRAINT uni_stores_name UNIQUE (name);


--
-- Name: users uni_users_email; Type: CONSTRAINT; Schema: public; Owner: admin
--

ALTER TABLE ONLY public.users
    ADD CONSTRAINT uni_users_email UNIQUE (email);


--
-- Name: users uni_users_username; Type: CONSTRAINT; Schema: public; Owner: admin
--

ALTER TABLE ONLY public.users
    ADD CONSTRAINT uni_users_username UNIQUE (username);


--
-- Name: users users_pkey; Type: CONSTRAINT; Schema: public; Owner: admin
--

ALTER TABLE ONLY public.users
    ADD CONSTRAINT users_pkey PRIMARY KEY (id);


--
-- Name: idx_order_items_deleted_at; Type: INDEX; Schema: public; Owner: admin
--

CREATE INDEX idx_order_items_deleted_at ON public.order_items USING btree (deleted_at);


--
-- Name: idx_orders_deleted_at; Type: INDEX; Schema: public; Owner: admin
--

CREATE INDEX idx_orders_deleted_at ON public.orders USING btree (deleted_at);


--
-- Name: idx_products_deleted_at; Type: INDEX; Schema: public; Owner: admin
--

CREATE INDEX idx_products_deleted_at ON public.products USING btree (deleted_at);


--
-- Name: idx_products_slug; Type: INDEX; Schema: public; Owner: admin
--

CREATE UNIQUE INDEX idx_products_slug ON public.products USING btree (slug);


--
-- Name: idx_stores_deleted_at; Type: INDEX; Schema: public; Owner: admin
--

CREATE INDEX idx_stores_deleted_at ON public.stores USING btree (deleted_at);


--
-- Name: idx_users_deleted_at; Type: INDEX; Schema: public; Owner: admin
--

CREATE INDEX idx_users_deleted_at ON public.users USING btree (deleted_at);


--
-- Name: order_items fk_order_items_product; Type: FK CONSTRAINT; Schema: public; Owner: admin
--

ALTER TABLE ONLY public.order_items
    ADD CONSTRAINT fk_order_items_product FOREIGN KEY (product_id) REFERENCES public.products(id);


--
-- Name: order_items fk_order_items_store; Type: FK CONSTRAINT; Schema: public; Owner: admin
--

ALTER TABLE ONLY public.order_items
    ADD CONSTRAINT fk_order_items_store FOREIGN KEY (store_id) REFERENCES public.stores(id);


--
-- Name: order_items fk_orders_order_items; Type: FK CONSTRAINT; Schema: public; Owner: admin
--

ALTER TABLE ONLY public.order_items
    ADD CONSTRAINT fk_orders_order_items FOREIGN KEY (order_id) REFERENCES public.orders(id);


--
-- Name: orders fk_orders_user; Type: FK CONSTRAINT; Schema: public; Owner: admin
--

ALTER TABLE ONLY public.orders
    ADD CONSTRAINT fk_orders_user FOREIGN KEY (user_id) REFERENCES public.users(id);


--
-- Name: products fk_products_store; Type: FK CONSTRAINT; Schema: public; Owner: admin
--

ALTER TABLE ONLY public.products
    ADD CONSTRAINT fk_products_store FOREIGN KEY (store_id) REFERENCES public.stores(id);


--
-- Name: stores fk_stores_user; Type: FK CONSTRAINT; Schema: public; Owner: admin
--

ALTER TABLE ONLY public.stores
    ADD CONSTRAINT fk_stores_user FOREIGN KEY (user_id) REFERENCES public.users(id);


--
-- PostgreSQL database dump complete
--

\unrestrict aJmBIO3L7VKA7TlJsdaeP6fiTN3gHuhh0Wi9Cy9jzvZPgxMtEHcCdGzpTOfSi6c

