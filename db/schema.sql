--
-- PostgreSQL database dump
--

\restrict Rm2Ko6u8biMQDBMZF9k8fajM5b1F7iVVNl4UkRQGxlgOlDTA8PCPjBeajHNuoNR

-- Dumped from database version 18.4
-- Dumped by pg_dump version 18.4

SET statement_timeout = 0;
SET lock_timeout = 0;
SET idle_in_transaction_session_timeout = 0;
SET transaction_timeout = 0;
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
-- Name: designs; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.designs (
    id integer NOT NULL,
    config jsonb NOT NULL,
    guest_token character varying(64) NOT NULL,
    created_at timestamp without time zone DEFAULT now()
);


ALTER TABLE public.designs OWNER TO postgres;

--
-- Name: designs_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.designs_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.designs_id_seq OWNER TO postgres;

--
-- Name: designs_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.designs_id_seq OWNED BY public.designs.id;


--
-- Name: orders; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.orders (
    id integer NOT NULL,
    design_id integer NOT NULL,
    guest_token character varying(64) NOT NULL,
    status character varying(20) DEFAULT 'pending'::character varying NOT NULL,
    amount numeric(10,2) NOT NULL,
    created_at timestamp without time zone DEFAULT now()
);


ALTER TABLE public.orders OWNER TO postgres;

--
-- Name: orders_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.orders_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.orders_id_seq OWNER TO postgres;

--
-- Name: orders_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.orders_id_seq OWNED BY public.orders.id;


--
-- Name: pdf_jobs; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.pdf_jobs (
    id integer NOT NULL,
    order_id integer NOT NULL,
    status character varying(20) DEFAULT 'pending'::character varying NOT NULL,
    config jsonb NOT NULL,
    file_url text,
    created_at timestamp without time zone DEFAULT now()
);


ALTER TABLE public.pdf_jobs OWNER TO postgres;

--
-- Name: pdf_jobs_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.pdf_jobs_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.pdf_jobs_id_seq OWNER TO postgres;

--
-- Name: pdf_jobs_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.pdf_jobs_id_seq OWNED BY public.pdf_jobs.id;


--
-- Name: designs id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.designs ALTER COLUMN id SET DEFAULT nextval('public.designs_id_seq'::regclass);


--
-- Name: orders id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.orders ALTER COLUMN id SET DEFAULT nextval('public.orders_id_seq'::regclass);


--
-- Name: pdf_jobs id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.pdf_jobs ALTER COLUMN id SET DEFAULT nextval('public.pdf_jobs_id_seq'::regclass);


--
-- Name: designs designs_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.designs
    ADD CONSTRAINT designs_pkey PRIMARY KEY (id);


--
-- Name: orders orders_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.orders
    ADD CONSTRAINT orders_pkey PRIMARY KEY (id);


--
-- Name: pdf_jobs pdf_jobs_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.pdf_jobs
    ADD CONSTRAINT pdf_jobs_pkey PRIMARY KEY (id);


--
-- Name: orders orders_design_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.orders
    ADD CONSTRAINT orders_design_id_fkey FOREIGN KEY (design_id) REFERENCES public.designs(id);


--
-- Name: pdf_jobs pdf_jobs_order_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.pdf_jobs
    ADD CONSTRAINT pdf_jobs_order_id_fkey FOREIGN KEY (order_id) REFERENCES public.orders(id);


--
-- PostgreSQL database dump complete
--

\unrestrict Rm2Ko6u8biMQDBMZF9k8fajM5b1F7iVVNl4UkRQGxlgOlDTA8PCPjBeajHNuoNR

