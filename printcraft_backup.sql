--
-- PostgreSQL database dump
--

\restrict x6IBz7RhGEoUbcRnMGwqGcdyuIyyIKCzYb74Q6PfNGU4ptToSGtheQRXk92corP

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
-- Name: keyboard_layouts; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.keyboard_layouts (
    id integer NOT NULL,
    alphabet_name character varying(50) NOT NULL,
    slot_count integer NOT NULL,
    keys jsonb NOT NULL,
    created_at timestamp without time zone DEFAULT now()
);


ALTER TABLE public.keyboard_layouts OWNER TO postgres;

--
-- Name: keyboard_layouts_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.keyboard_layouts_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.keyboard_layouts_id_seq OWNER TO postgres;

--
-- Name: keyboard_layouts_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.keyboard_layouts_id_seq OWNED BY public.keyboard_layouts.id;


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
-- Name: keyboard_layouts id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.keyboard_layouts ALTER COLUMN id SET DEFAULT nextval('public.keyboard_layouts_id_seq'::regclass);


--
-- Name: orders id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.orders ALTER COLUMN id SET DEFAULT nextval('public.orders_id_seq'::regclass);


--
-- Name: pdf_jobs id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.pdf_jobs ALTER COLUMN id SET DEFAULT nextval('public.pdf_jobs_id_seq'::regclass);


--
-- Data for Name: designs; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.designs (id, config, guest_token, created_at) FROM stdin;
2	{"color": "green", "alphabet": "hiragana"}	xyz789	2026-08-02 19:33:14.345569
3	{"color": "purple", "alphabet": "arabic"}	test999	2026-08-03 23:59:05.034836
4	{"key_mode": "black", "primary_alphabet": "hiragana", "secondary_alphabet": null}	t1	2026-08-23 04:10:12.791931
5	{"key_mode": "white", "primary_alphabet": "azerty_fr", "secondary_alphabet": "hiragana"}	t2	2026-08-23 04:10:12.791931
\.


--
-- Data for Name: keyboard_layouts; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.keyboard_layouts (id, alphabet_name, slot_count, keys, created_at) FROM stdin;
3	hiragana	1	[{"base": "あ", "key_code": "A"}, {"base": "い", "key_code": "B"}, {"base": "う", "key_code": "C"}, {"base": "え", "key_code": "D"}]	2026-08-22 18:54:56.441044
4	azerty_fr	4	[{"base": "²", "altgr": null, "shift": null, "key_code": "k0", "altgr_shift": null}, {"base": "&", "altgr": null, "shift": "1", "key_code": "k1", "altgr_shift": null}, {"base": "é", "altgr": "~", "shift": "2", "key_code": "k2", "altgr_shift": null}, {"base": "\\"", "altgr": "#", "shift": "3", "key_code": "k3", "altgr_shift": null}, {"base": "'", "altgr": "{", "shift": "4", "key_code": "k4", "altgr_shift": null}, {"base": "(", "altgr": "[", "shift": "5", "key_code": "k5", "altgr_shift": null}, {"base": "-", "altgr": "|", "shift": "6", "key_code": "k6", "altgr_shift": null}, {"base": "è", "altgr": null, "shift": "7", "key_code": "k7", "altgr_shift": null}, {"base": "_", "altgr": "\\\\", "shift": "8", "key_code": "k8", "altgr_shift": null}, {"base": "ç", "altgr": "^", "shift": "9", "key_code": "k9", "altgr_shift": null}, {"base": "à", "altgr": "@", "shift": "0", "key_code": "k10", "altgr_shift": null}, {"base": ")", "altgr": "]", "shift": "°", "key_code": "k11", "altgr_shift": null}, {"base": "=", "altgr": "}", "shift": "+", "key_code": "k12", "altgr_shift": null}, {"base": "^", "altgr": null, "shift": "¨", "key_code": "k13", "altgr_shift": null}, {"base": "$", "altgr": "¤", "shift": "£", "key_code": "k14", "altgr_shift": null}, {"base": "*", "altgr": null, "shift": "µ", "key_code": "k15", "altgr_shift": null}, {"base": "ù", "altgr": null, "shift": "%", "key_code": "k16", "altgr_shift": null}, {"base": "!", "altgr": null, "shift": "§", "key_code": "k17", "altgr_shift": null}, {"base": ":", "altgr": null, "shift": "/", "key_code": "k18", "altgr_shift": null}, {"base": ";", "altgr": null, "shift": ".", "key_code": "k19", "altgr_shift": null}, {"base": ",", "altgr": null, "shift": "?", "key_code": "k20", "altgr_shift": null}, {"base": "E", "altgr": "€", "shift": null, "key_code": "k21", "altgr_shift": null}, {"base": "A", "altgr": null, "shift": null, "key_code": "k22", "altgr_shift": null}, {"base": "Z", "altgr": null, "shift": null, "key_code": "k23", "altgr_shift": null}, {"base": "R", "altgr": null, "shift": null, "key_code": "k24", "altgr_shift": null}, {"base": "T", "altgr": null, "shift": null, "key_code": "k25", "altgr_shift": null}, {"base": "Y", "altgr": null, "shift": null, "key_code": "k26", "altgr_shift": null}, {"base": "U", "altgr": null, "shift": null, "key_code": "k27", "altgr_shift": null}, {"base": "I", "altgr": null, "shift": null, "key_code": "k28", "altgr_shift": null}, {"base": "O", "altgr": null, "shift": null, "key_code": "k29", "altgr_shift": null}, {"base": "P", "altgr": null, "shift": null, "key_code": "k30", "altgr_shift": null}, {"base": "Q", "altgr": null, "shift": null, "key_code": "k31", "altgr_shift": null}, {"base": "S", "altgr": null, "shift": null, "key_code": "k32", "altgr_shift": null}, {"base": "D", "altgr": null, "shift": null, "key_code": "k33", "altgr_shift": null}, {"base": "F", "altgr": null, "shift": null, "key_code": "k34", "altgr_shift": null}, {"base": "G", "altgr": null, "shift": null, "key_code": "k35", "altgr_shift": null}, {"base": "H", "altgr": null, "shift": null, "key_code": "k36", "altgr_shift": null}, {"base": "J", "altgr": null, "shift": null, "key_code": "k37", "altgr_shift": null}, {"base": "K", "altgr": null, "shift": null, "key_code": "k38", "altgr_shift": null}, {"base": "L", "altgr": null, "shift": null, "key_code": "k39", "altgr_shift": null}, {"base": "M", "altgr": null, "shift": null, "key_code": "k40", "altgr_shift": null}, {"base": "W", "altgr": null, "shift": null, "key_code": "k41", "altgr_shift": null}, {"base": "X", "altgr": null, "shift": null, "key_code": "k42", "altgr_shift": null}, {"base": "C", "altgr": null, "shift": null, "key_code": "k43", "altgr_shift": null}, {"base": "V", "altgr": null, "shift": null, "key_code": "k44", "altgr_shift": null}, {"base": "B", "altgr": null, "shift": null, "key_code": "k45", "altgr_shift": null}, {"base": "N", "altgr": null, "shift": null, "key_code": "k46", "altgr_shift": null}]	2026-08-22 18:54:56.441044
\.


--
-- Data for Name: orders; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.orders (id, design_id, guest_token, status, amount, created_at) FROM stdin;
1	2	xyz789	paid	5.99	2026-08-02 19:33:37.824585
2	3	test999	paid	4.99	2026-08-03 23:59:31.69852
\.


--
-- Data for Name: pdf_jobs; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.pdf_jobs (id, order_id, status, config, file_url, created_at) FROM stdin;
1	1	pending	{"color": "green", "alphabet": "hiragana"}	\N	2026-08-02 19:33:37.824585
2	2	pending	{"color": "purple", "alphabet": "arabic"}	\N	2026-08-03 23:59:31.69852
\.


--
-- Name: designs_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.designs_id_seq', 5, true);


--
-- Name: keyboard_layouts_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.keyboard_layouts_id_seq', 4, true);


--
-- Name: orders_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.orders_id_seq', 2, true);


--
-- Name: pdf_jobs_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.pdf_jobs_id_seq', 2, true);


--
-- Name: designs designs_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.designs
    ADD CONSTRAINT designs_pkey PRIMARY KEY (id);


--
-- Name: keyboard_layouts keyboard_layouts_alphabet_name_key; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.keyboard_layouts
    ADD CONSTRAINT keyboard_layouts_alphabet_name_key UNIQUE (alphabet_name);


--
-- Name: keyboard_layouts keyboard_layouts_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.keyboard_layouts
    ADD CONSTRAINT keyboard_layouts_pkey PRIMARY KEY (id);


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

\unrestrict x6IBz7RhGEoUbcRnMGwqGcdyuIyyIKCzYb74Q6PfNGU4ptToSGtheQRXk92corP

