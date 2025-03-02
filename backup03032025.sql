--
-- PostgreSQL database dump
--

-- Dumped from database version 15.12 (Debian 15.12-1.pgdg120+1)
-- Dumped by pg_dump version 15.12 (Debian 15.12-1.pgdg120+1)

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
-- Name: datastamps; Type: TABLE; Schema: public; Owner: tanki_enjoyer
--

CREATE TABLE public.datastamps (
    datastamp_id integer NOT NULL,
    user_id integer NOT NULL,
    created_at timestamp without time zone NOT NULL,
    rank smallint,
    kills integer,
    deaths integer,
    cry integer
);


ALTER TABLE public.datastamps OWNER TO tanki_enjoyer;

--
-- Name: datastamps_datastamp_id_seq; Type: SEQUENCE; Schema: public; Owner: tanki_enjoyer
--

CREATE SEQUENCE public.datastamps_datastamp_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.datastamps_datastamp_id_seq OWNER TO tanki_enjoyer;

--
-- Name: datastamps_datastamp_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: tanki_enjoyer
--

ALTER SEQUENCE public.datastamps_datastamp_id_seq OWNED BY public.datastamps.datastamp_id;


--
-- Name: gear_stats; Type: TABLE; Schema: public; Owner: tanki_enjoyer
--

CREATE TABLE public.gear_stats (
    datastamp_id integer NOT NULL,
    gear_key smallint NOT NULL,
    score_earned integer NOT NULL,
    seconds_played integer NOT NULL
);


ALTER TABLE public.gear_stats OWNER TO tanki_enjoyer;

--
-- Name: users; Type: TABLE; Schema: public; Owner: tanki_enjoyer
--

CREATE TABLE public.users (
    user_id integer NOT NULL,
    name character varying(128) NOT NULL
);


ALTER TABLE public.users OWNER TO tanki_enjoyer;

--
-- Name: users_user_id_seq; Type: SEQUENCE; Schema: public; Owner: tanki_enjoyer
--

CREATE SEQUENCE public.users_user_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.users_user_id_seq OWNER TO tanki_enjoyer;

--
-- Name: users_user_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: tanki_enjoyer
--

ALTER SEQUENCE public.users_user_id_seq OWNED BY public.users.user_id;


--
-- Name: datastamps datastamp_id; Type: DEFAULT; Schema: public; Owner: tanki_enjoyer
--

ALTER TABLE ONLY public.datastamps ALTER COLUMN datastamp_id SET DEFAULT nextval('public.datastamps_datastamp_id_seq'::regclass);


--
-- Name: users user_id; Type: DEFAULT; Schema: public; Owner: tanki_enjoyer
--

ALTER TABLE ONLY public.users ALTER COLUMN user_id SET DEFAULT nextval('public.users_user_id_seq'::regclass);


--
-- Data for Name: datastamps; Type: TABLE DATA; Schema: public; Owner: tanki_enjoyer
--

COPY public.datastamps (datastamp_id, user_id, created_at, rank, kills, deaths, cry) FROM stdin;
1	1	2025-03-02 14:00:00	0	0	0	0
2	5	2025-03-02 14:00:00	129	215064	130081	15821168
3	2	2025-03-02 14:00:00	98	215605	107425	12000302
4	4	2025-03-02 14:00:00	136	226092	134549	20807270
5	1	2025-03-02 15:00:00	0	0	0	0
6	4	2025-03-02 15:00:00	136	226092	134549	20807270
7	2	2025-03-02 15:00:00	98	215605	107425	12000302
8	5	2025-03-02 15:00:00	129	215064	130081	15821168
9	1	2025-03-02 16:00:00	0	0	0	0
10	5	2025-03-02 16:00:00	129	215064	130081	15821168
11	2	2025-03-02 16:00:00	98	215626	107430	12000302
12	4	2025-03-02 16:00:00	136	226092	134549	20807270
13	1	2025-03-02 18:00:00	0	0	0	0
14	2	2025-03-02 18:00:00	99	215711	107478	12004405
15	5	2025-03-02 18:00:00	129	215064	130081	15821168
16	4	2025-03-02 18:00:00	136	226092	134549	20807270
17	1	2025-03-02 21:00:00	0	0	0	0
18	2	2025-03-02 21:00:00	99	215711	107478	12004405
19	5	2025-03-02 21:00:00	129	215064	130081	15821168
20	4	2025-03-02 21:00:00	136	226092	134549	20807270
\.


--
-- Data for Name: gear_stats; Type: TABLE DATA; Schema: public; Owner: tanki_enjoyer
--

COPY public.gear_stats (datastamp_id, gear_key, score_earned, seconds_played) FROM stdin;
6	1	1542708	568329
6	2	3265574	1630854
6	3	4520375	1621483
6	4	1422186	481890
6	5	1596356	574636
6	7	1337506	440168
6	9	282350	162167
6	10	1057570	374846
6	0	1780649	1092089
6	6	3967196	1803620
6	8	499842	224989
6	109	537417	194013
6	110	972825	406661
6	102	1001853	434686
6	106	933000	452712
6	112	469421	212086
6	115	1803133	809613
6	111	1447444	776141
6	113	569067	192470
6	101	815082	387932
6	108	1691524	470789
6	104	957182	305698
6	105	665221	413841
6	107	780507	482441
6	114	559304	228306
6	100	482251	548008
6	103	7746733	2700524
7	7	463521	140303
7	8	89	1739
7	9	54474	18999
7	0	11704467	4794882
7	1	1023391	440041
7	3	579579	206061
7	5	6042	2625
7	10	109444	39526
7	2	73436	23857
7	4	1034061	378866
7	106	135110	44712
7	110	375810	151235
7	115	1370658	601436
7	100	373759	131781
7	102	243869	91357
7	103	785901	245425
7	107	34527	18884
7	114	35486	14121
7	105	2864	2761
7	109	492938	147273
7	111	357959	126493
7	112	770968	487542
7	113	5144066	2311217
7	101	276911	75628
7	104	925847	317439
7	108	3720631	1279286
8	8	1500641	617682
8	10	172744	47796
8	1	71113	24569
8	2	3682122	1171665
8	4	1418196	647068
8	5	4344037	1371257
8	6	1022152	337881
8	0	1012592	482079
8	3	3971960	1180858
8	7	1488997	343451
8	9	69283	26667
8	107	370249	119058
8	109	466162	116514
8	110	2930951	972297
8	112	1220836	608090
8	114	1311046	365265
8	102	102414	49508
8	103	340131	85811
8	105	144764	41400
8	101	823758	256129
8	104	2439468	740667
8	108	2673750	923735
8	111	294612	156149
8	113	318888	138344
8	100	1382881	313475
8	106	1873957	645225
8	115	2094252	729551
11	0	11704547	4794908
11	3	580419	206462
11	113	5144986	2311643
14	0	11707307	4796587
14	1	1027421	441543
14	112	773028	488848
14	113	5145686	2312017
14	110	379840	152737
\.


--
-- Data for Name: users; Type: TABLE DATA; Schema: public; Owner: tanki_enjoyer
--

COPY public.users (user_id, name) FROM stdin;
1	Zlodeath
2	silly
4	Isshiki
5	elastic
\.


--
-- Name: datastamps_datastamp_id_seq; Type: SEQUENCE SET; Schema: public; Owner: tanki_enjoyer
--

SELECT pg_catalog.setval('public.datastamps_datastamp_id_seq', 20, true);


--
-- Name: users_user_id_seq; Type: SEQUENCE SET; Schema: public; Owner: tanki_enjoyer
--

SELECT pg_catalog.setval('public.users_user_id_seq', 5, true);


--
-- Name: datastamps datastamps_pkey; Type: CONSTRAINT; Schema: public; Owner: tanki_enjoyer
--

ALTER TABLE ONLY public.datastamps
    ADD CONSTRAINT datastamps_pkey PRIMARY KEY (datastamp_id);


--
-- Name: gear_stats gear_stats_pkey; Type: CONSTRAINT; Schema: public; Owner: tanki_enjoyer
--

ALTER TABLE ONLY public.gear_stats
    ADD CONSTRAINT gear_stats_pkey PRIMARY KEY (datastamp_id, gear_key);


--
-- Name: users users_name_key; Type: CONSTRAINT; Schema: public; Owner: tanki_enjoyer
--

ALTER TABLE ONLY public.users
    ADD CONSTRAINT users_name_key UNIQUE (name);


--
-- Name: users users_pkey; Type: CONSTRAINT; Schema: public; Owner: tanki_enjoyer
--

ALTER TABLE ONLY public.users
    ADD CONSTRAINT users_pkey PRIMARY KEY (user_id);


--
-- Name: idx_datastamps_created_at; Type: INDEX; Schema: public; Owner: tanki_enjoyer
--

CREATE INDEX idx_datastamps_created_at ON public.datastamps USING btree (created_at);


--
-- Name: datastamps datastamps_user_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: tanki_enjoyer
--

ALTER TABLE ONLY public.datastamps
    ADD CONSTRAINT datastamps_user_id_fkey FOREIGN KEY (user_id) REFERENCES public.users(user_id) ON DELETE CASCADE;


--
-- Name: gear_stats gear_stats_datastamp_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: tanki_enjoyer
--

ALTER TABLE ONLY public.gear_stats
    ADD CONSTRAINT gear_stats_datastamp_id_fkey FOREIGN KEY (datastamp_id) REFERENCES public.datastamps(datastamp_id) ON DELETE CASCADE;


--
-- PostgreSQL database dump complete
--

