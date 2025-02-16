package contenttype

import "strings"

// List of ISO 639 set 1 language codes
var languageSet1 = map[string]string{
	"aa": "aar",
	"ab": "abk",
	"ae": "ave",
	"af": "afr",
	"ak": "aka",
	"am": "amh",
	"an": "arg",
	"ar": "ara",
	"as": "asm",
	"av": "ava",
	"ay": "aym",
	"az": "aze",
	"ba": "bak",
	"be": "bel",
	"bg": "bul",
	"bh": "bih",
	"bi": "bis",
	"bm": "bam",
	"bn": "ben",
	"bo": "bod",
	"br": "bre",
	"bs": "bos",
	"ca": "cat",
	"ce": "che",
	"ch": "cha",
	"co": "cos",
	"cr": "cre",
	"cs": "ces",
	"cu": "chu",
	"cv": "chv",
	"cy": "cym",
	"da": "dan",
	"de": "deu",
	"dv": "div",
	"dz": "dzo",
	"ee": "ewe",
	"el": "ell",
	"en": "eng",
	"eo": "epo",
	"es": "spa",
	"et": "est",
	"eu": "eus",
	"fa": "fas",
	"ff": "ful",
	"fi": "fin",
	"fj": "fij",
	"fo": "fao",
	"fr": "fra",
	"fy": "fry",
	"ga": "gle",
	"gd": "gla",
	"gl": "glg",
	"gn": "grn",
	"gu": "guj",
	"gv": "glv",
	"ha": "hau",
	"he": "heb",
	"hi": "hin",
	"ho": "hmo",
	"hr": "hrv",
	"ht": "hat",
	"hu": "hun",
	"hy": "hye",
	"hz": "her",
	"ia": "ina",
	"id": "ind",
	"ie": "ile",
	"ig": "ibo",
	"ii": "iii",
	"ik": "ipk",
	"io": "ido",
	"is": "isl",
	"it": "ita",
	"iu": "iku",
	"ja": "jpn",
	"jv": "jav",
	"ka": "kat",
	"kg": "kon",
	"ki": "kik",
	"kj": "kua",
	"kk": "kaz",
	"kl": "kal",
	"km": "khm",
	"kn": "kan",
	"ko": "kor",
	"kr": "kau",
	"ks": "kas",
	"ku": "kur",
	"kv": "kom",
	"kw": "cor",
	"ky": "kir",
	"la": "lat",
	"lb": "ltz",
	"lg": "lug",
	"li": "lim",
	"ln": "lin",
	"lo": "lao",
	"lt": "lit",
	"lu": "lub",
	"lv": "lav",
	"mg": "mlg",
	"mh": "mah",
	"mi": "mri",
	"mk": "mkd",
	"ml": "mal",
	"mn": "mon",
	"mr": "mar",
	"ms": "msa",
	"mt": "mlt",
	"my": "mya",
	"na": "nau",
	"nb": "nob",
	"nd": "nde",
	"ne": "nep",
	"ng": "ndo",
	"nl": "nld",
	"nn": "nno",
	"no": "nor",
	"nr": "nbl",
	"nv": "nav",
	"ny": "nya",
	"oc": "oci",
	"oj": "oji",
	"om": "orm",
	"or": "ori",
	"os": "oss",
	"pa": "pan",
	"pi": "pli",
	"pl": "pol",
	"ps": "pus",
	"pt": "por",
	"qu": "que",
	"rm": "roh",
	"rn": "run",
	"ro": "ron",
	"ru": "rus",
	"rw": "kin",
	"sa": "san",
	"sc": "srd",
	"sd": "snd",
	"se": "sme",
	"sg": "sag",
	"si": "sin",
	"sk": "slk",
	"sl": "slv",
	"sm": "smo",
	"sn": "sna",
	"so": "som",
	"sq": "sqi",
	"sr": "srp",
	"ss": "ssw",
	"st": "sot",
	"su": "sun",
	"sv": "swe",
	"sw": "swa",
	"ta": "tam",
	"te": "tel",
	"tg": "tgk",
	"th": "tha",
	"ti": "tir",
	"tk": "tuk",
	"tl": "tgl",
	"tn": "tsn",
	"to": "ton",
	"tr": "tur",
	"ts": "tso",
	"tt": "tat",
	"tw": "twi",
	"ty": "tah",
	"ug": "uig",
	"uk": "ukr",
	"ur": "urd",
	"uz": "uzb",
	"ve": "ven",
	"vi": "vie",
	"vo": "vol",
	"wa": "wln",
	"wo": "wol",
	"xh": "xho",
	"yi": "yid",
	"yo": "yor",
	"za": "zha",
	"zh": "zho",
	"zu": "zul",
}

// List of ISO 639 set 2 language codes
var languageSet2 = map[string]string{
	"aar": "aar",
	"abk": "abk",
	"ace": "ace",
	"ach": "ach",
	"ada": "ada",
	"ady": "ady",
	"afa": "afa",
	"afh": "afh",
	"afr": "afr",
	"ain": "ain",
	"aka": "aka",
	"akk": "akk",
	"alb": "sqi",
	"ale": "ale",
	"alg": "alg",
	"alt": "alt",
	"amh": "amh",
	"ang": "ang",
	"anp": "anp",
	"apa": "apa",
	"ara": "ara",
	"arc": "arc",
	"arg": "arg",
	"arm": "hye",
	"arn": "arn",
	"arp": "arp",
	"art": "art",
	"arw": "arw",
	"asm": "asm",
	"ast": "ast",
	"ath": "ath",
	"aus": "aus",
	"ava": "ava",
	"ave": "ave",
	"awa": "awa",
	"aym": "aym",
	"aze": "aze",
	"bad": "bad",
	"bai": "bai",
	"bak": "bak",
	"bal": "bal",
	"bam": "bam",
	"ban": "ban",
	"baq": "eus",
	"bas": "bas",
	"bat": "bat",
	"bej": "bej",
	"bel": "bel",
	"bem": "bem",
	"ben": "ben",
	"ber": "ber",
	"bho": "bho",
	"bih": "bih",
	"bik": "bik",
	"bin": "bin",
	"bis": "bis",
	"bla": "bla",
	"bnt": "bnt",
	"bod": "bod",
	"bos": "bos",
	"bra": "bra",
	"bre": "bre",
	"btk": "btk",
	"bua": "bua",
	"bug": "bug",
	"bul": "bul",
	"bur": "mya",
	"byn": "byn",
	"cad": "cad",
	"cai": "cai",
	"car": "car",
	"cat": "cat",
	"cau": "cau",
	"ceb": "ceb",
	"cel": "cel",
	"ces": "ces",
	"cha": "cha",
	"chb": "chb",
	"che": "che",
	"chg": "chg",
	"chi": "zho",
	"chk": "chk",
	"chm": "chm",
	"chn": "chn",
	"cho": "cho",
	"chp": "chp",
	"chr": "chr",
	"chu": "chu",
	"chv": "chv",
	"chy": "chy",
	"cmc": "cmc",
	"cop": "cop",
	"cor": "cor",
	"cos": "cos",
	"cpe": "cpe",
	"cpf": "cpf",
	"cpp": "cpp",
	"cre": "cre",
	"crh": "crh",
	"crp": "crp",
	"csb": "csb",
	"cus": "cus",
	"cym": "cym",
	"cze": "ces",
	"dak": "dak",
	"dan": "dan",
	"dar": "dar",
	"day": "day",
	"del": "del",
	"den": "den",
	"deu": "deu",
	"dgr": "dgr",
	"din": "din",
	"div": "div",
	"doi": "doi",
	"dra": "dra",
	"dsb": "dsb",
	"dua": "dua",
	"dum": "dum",
	"dut": "nld",
	"dyu": "dyu",
	"dzo": "dzo",
	"efi": "efi",
	"egy": "egy",
	"eka": "eka",
	"ell": "ell",
	"elx": "elx",
	"eng": "eng",
	"enm": "enm",
	"epo": "epo",
	"est": "est",
	"eus": "eus",
	"ewe": "ewe",
	"ewo": "ewo",
	"fan": "fan",
	"fao": "fao",
	"fas": "fas",
	"fat": "fat",
	"fij": "fij",
	"fil": "fil",
	"fin": "fin",
	"fiu": "fiu",
	"fon": "fon",
	"fra": "fra",
	"fre": "fra",
	"frm": "frm",
	"fro": "fro",
	"frr": "frr",
	"frs": "frs",
	"fry": "fry",
	"ful": "ful",
	"fur": "fur",
	"gaa": "gaa",
	"gay": "gay",
	"gba": "gba",
	"gem": "gem",
	"geo": "kat",
	"ger": "deu",
	"gez": "gez",
	"gil": "gil",
	"gla": "gla",
	"gle": "gle",
	"glg": "glg",
	"glv": "glv",
	"gmh": "gmh",
	"goh": "goh",
	"gon": "gon",
	"gor": "gor",
	"got": "got",
	"grb": "grb",
	"grc": "grc",
	"gre": "ell",
	"grn": "grn",
	"gsw": "gsw",
	"guj": "guj",
	"gwi": "gwi",
	"hai": "hai",
	"hat": "hat",
	"hau": "hau",
	"haw": "haw",
	"heb": "heb",
	"her": "her",
	"hil": "hil",
	"him": "him",
	"hin": "hin",
	"hit": "hit",
	"hmn": "hmn",
	"hmo": "hmo",
	"hrv": "hrv",
	"hsb": "hsb",
	"hun": "hun",
	"hup": "hup",
	"hye": "hye",
	"iba": "iba",
	"ibo": "ibo",
	"ice": "isl",
	"ido": "ido",
	"iii": "iii",
	"ijo": "ijo",
	"iku": "iku",
	"ile": "ile",
	"ilo": "ilo",
	"ina": "ina",
	"inc": "inc",
	"ind": "ind",
	"ine": "ine",
	"inh": "inh",
	"ipk": "ipk",
	"ira": "ira",
	"iro": "iro",
	"isl": "isl",
	"ita": "ita",
	"jav": "jav",
	"jbo": "jbo",
	"jpn": "jpn",
	"jpr": "jpr",
	"jrb": "jrb",
	"kaa": "kaa",
	"kab": "kab",
	"kac": "kac",
	"kal": "kal",
	"kam": "kam",
	"kan": "kan",
	"kar": "kar",
	"kas": "kas",
	"kat": "kat",
	"kau": "kau",
	"kaw": "kaw",
	"kaz": "kaz",
	"kbd": "kbd",
	"kha": "kha",
	"khi": "khi",
	"khm": "khm",
	"kho": "kho",
	"kik": "kik",
	"kin": "kin",
	"kir": "kir",
	"kmb": "kmb",
	"kok": "kok",
	"kom": "kom",
	"kon": "kon",
	"kor": "kor",
	"kos": "kos",
	"kpe": "kpe",
	"krc": "krc",
	"krl": "krl",
	"kro": "kro",
	"kru": "kru",
	"kua": "kua",
	"kum": "kum",
	"kur": "kur",
	"kut": "kut",
	"lad": "lad",
	"lah": "lah",
	"lam": "lam",
	"lao": "lao",
	"lat": "lat",
	"lav": "lav",
	"lez": "lez",
	"lim": "lim",
	"lin": "lin",
	"lit": "lit",
	"lol": "lol",
	"loz": "loz",
	"ltz": "ltz",
	"lua": "lua",
	"lub": "lub",
	"lug": "lug",
	"lui": "lui",
	"lun": "lun",
	"luo": "luo",
	"lus": "lus",
	"mac": "mkd",
	"mad": "mad",
	"mag": "mag",
	"mah": "mah",
	"mai": "mai",
	"mak": "mak",
	"mal": "mal",
	"man": "man",
	"mao": "mri",
	"map": "map",
	"mar": "mar",
	"mas": "mas",
	"may": "msa",
	"mdf": "mdf",
	"mdr": "mdr",
	"men": "men",
	"mga": "mga",
	"mic": "mic",
	"min": "min",
	"mis": "mis",
	"mkd": "mkd",
	"mkh": "mkh",
	"mlg": "mlg",
	"mlt": "mlt",
	"mnc": "mnc",
	"mni": "mni",
	"mno": "mno",
	"moh": "moh",
	"mon": "mon",
	"mos": "mos",
	"mri": "mri",
	"msa": "msa",
	"mul": "mul",
	"mun": "mun",
	"mus": "mus",
	"mwl": "mwl",
	"mwr": "mwr",
	"mya": "mya",
	"myn": "myn",
	"myv": "myv",
	"nah": "nah",
	"nai": "nai",
	"nap": "nap",
	"nau": "nau",
	"nav": "nav",
	"nbl": "nbl",
	"nde": "nde",
	"ndo": "ndo",
	"nds": "nds",
	"nep": "nep",
	"new": "new",
	"nia": "nia",
	"nic": "nic",
	"niu": "niu",
	"nld": "nld",
	"nno": "nno",
	"nob": "nob",
	"nog": "nog",
	"non": "non",
	"nor": "nor",
	"nqo": "nqo",
	"nso": "nso",
	"nub": "nub",
	"nwc": "nwc",
	"nya": "nya",
	"nym": "nym",
	"nyn": "nyn",
	"nyo": "nyo",
	"nzi": "nzi",
	"oci": "oci",
	"oji": "oji",
	"ori": "ori",
	"orm": "orm",
	"osa": "osa",
	"oss": "oss",
	"ota": "ota",
	"oto": "oto",
	"paa": "paa",
	"pag": "pag",
	"pal": "pal",
	"pam": "pam",
	"pan": "pan",
	"pap": "pap",
	"pau": "pau",
	"peo": "peo",
	"per": "fas",
	"phi": "phi",
	"phn": "phn",
	"pli": "pli",
	"pol": "pol",
	"pon": "pon",
	"por": "por",
	"pra": "pra",
	"pro": "pro",
	"pus": "pus",
	"que": "que",
	"raj": "raj",
	"rap": "rap",
	"rar": "rar",
	"roa": "roa",
	"roh": "roh",
	"rom": "rom",
	"ron": "ron",
	"rum": "ron",
	"run": "run",
	"rup": "rup",
	"rus": "rus",
	"sad": "sad",
	"sag": "sag",
	"sah": "sah",
	"sai": "sai",
	"sal": "sal",
	"sam": "sam",
	"san": "san",
	"sas": "sas",
	"sat": "sat",
	"scn": "scn",
	"sco": "sco",
	"sel": "sel",
	"sem": "sem",
	"sga": "sga",
	"sgn": "sgn",
	"shn": "shn",
	"sid": "sid",
	"sin": "sin",
	"sio": "sio",
	"sit": "sit",
	"sla": "sla",
	"slk": "slk",
	"slo": "slk",
	"slv": "slv",
	"sma": "sma",
	"sme": "sme",
	"smi": "smi",
	"smj": "smj",
	"smn": "smn",
	"smo": "smo",
	"sms": "sms",
	"sna": "sna",
	"snd": "snd",
	"snk": "snk",
	"sog": "sog",
	"som": "som",
	"son": "son",
	"sot": "sot",
	"spa": "spa",
	"sqi": "sqi",
	"srd": "srd",
	"srn": "srn",
	"srp": "srp",
	"srr": "srr",
	"ssa": "ssa",
	"ssw": "ssw",
	"suk": "suk",
	"sun": "sun",
	"sus": "sus",
	"sux": "sux",
	"swa": "swa",
	"swe": "swe",
	"syc": "syc",
	"syr": "syr",
	"tah": "tah",
	"tai": "tai",
	"tam": "tam",
	"tat": "tat",
	"tel": "tel",
	"tem": "tem",
	"ter": "ter",
	"tet": "tet",
	"tgk": "tgk",
	"tgl": "tgl",
	"tha": "tha",
	"tib": "bod",
	"tig": "tig",
	"tir": "tir",
	"tiv": "tiv",
	"tkl": "tkl",
	"tlh": "tlh",
	"tli": "tli",
	"tmh": "tmh",
	"tog": "tog",
	"ton": "ton",
	"tpi": "tpi",
	"tsi": "tsi",
	"tsn": "tsn",
	"tso": "tso",
	"tuk": "tuk",
	"tum": "tum",
	"tup": "tup",
	"tur": "tur",
	"tut": "tut",
	"tvl": "tvl",
	"twi": "twi",
	"tyv": "tyv",
	"udm": "udm",
	"uga": "uga",
	"uig": "uig",
	"ukr": "ukr",
	"umb": "umb",
	"und": "und",
	"urd": "urd",
	"uzb": "uzb",
	"vai": "vai",
	"ven": "ven",
	"vie": "vie",
	"vol": "vol",
	"vot": "vot",
	"wak": "wak",
	"wal": "wal",
	"war": "war",
	"was": "was",
	"wel": "cym",
	"wen": "wen",
	"wln": "wln",
	"wol": "wol",
	"xal": "xal",
	"xho": "xho",
	"yao": "yao",
	"yap": "yap",
	"yid": "yid",
	"yor": "yor",
	"ypk": "ypk",
	"zap": "zap",
	"zbl": "zbl",
	"zen": "zen",
	"zgh": "zgh",
	"zha": "zha",
	"zho": "zho",
	"znd": "znd",
	"zul": "zul",
	"zun": "zun",
	"zxx": "zxx",
	"zza": "zza",
}

// List of ISO 15924 scripts
var scripts = map[string]int{
	"adlm": 166,
	"afak": 439,
	"aghb": 239,
	"ahom": 338,
	"arab": 160,
	"aran": 161,
	"armi": 124,
	"armn": 230,
	"avst": 134,
	"bali": 360,
	"bamu": 435,
	"bass": 259,
	"batk": 365,
	"beng": 325,
	"berf": 258,
	"bhks": 334,
	"blis": 550,
	"bopo": 285,
	"brah": 300,
	"brai": 570,
	"bugi": 367,
	"buhd": 372,
	"cakm": 349,
	"cans": 440,
	"cari": 201,
	"cham": 358,
	"cher": 445,
	"chis": 298,
	"chrs": 109,
	"cirt": 291,
	"copt": 204,
	"cpmn": 402,
	"cprt": 403,
	"cyrl": 220,
	"cyrs": 221,
	"deva": 315,
	"diak": 342,
	"dogr": 328,
	"dsrt": 250,
	"dupl": 755,
	"egyd": 70,
	"egyh": 60,
	"egyp": 50,
	"elba": 226,
	"elym": 128,
	"ethi": 430,
	"gara": 164,
	"geok": 241,
	"geor": 240,
	"glag": 225,
	"gong": 312,
	"gonm": 313,
	"goth": 206,
	"gran": 343,
	"grek": 200,
	"gujr": 320,
	"gukh": 397,
	"guru": 310,
	"hanb": 503,
	"hang": 286,
	"hani": 500,
	"hano": 371,
	"hans": 501,
	"hant": 502,
	"hatr": 127,
	"hebr": 125,
	"hira": 410,
	"hluw": 80,
	"hmng": 450,
	"hmnp": 451,
	"hrkt": 412,
	"hung": 176,
	"inds": 610,
	"ital": 210,
	"jamo": 284,
	"java": 361,
	"jpan": 413,
	"jurc": 510,
	"kali": 357,
	"kana": 411,
	"kawi": 368,
	"khar": 305,
	"khmr": 355,
	"khoj": 322,
	"kitl": 505,
	"kits": 288,
	"knda": 345,
	"kore": 287,
	"kpel": 436,
	"krai": 396,
	"kthi": 317,
	"lana": 351,
	"laoo": 356,
	"latf": 217,
	"latg": 216,
	"latn": 215,
	"leke": 364,
	"lepc": 335,
	"limb": 336,
	"lina": 400,
	"linb": 401,
	"lisu": 399,
	"loma": 437,
	"lyci": 202,
	"lydi": 116,
	"mahj": 314,
	"maka": 366,
	"mand": 140,
	"mani": 139,
	"marc": 332,
	"maya": 90,
	"medf": 265,
	"mend": 438,
	"merc": 101,
	"mero": 100,
	"mlym": 347,
	"modi": 324,
	"mong": 145,
	"moon": 218,
	"mroo": 264,
	"mtei": 337,
	"mult": 323,
	"mymr": 350,
	"nagm": 295,
	"nand": 311,
	"narb": 106,
	"nbat": 159,
	"newa": 333,
	"nkdb": 85,
	"nkgb": 420,
	"nkoo": 165,
	"nshu": 499,
	"ogam": 212,
	"olck": 261,
	"onao": 296,
	"orkh": 175,
	"orya": 327,
	"osge": 219,
	"osma": 260,
	"ougr": 143,
	"palm": 126,
	"pauc": 263,
	"pcun": 15,
	"pelm": 16,
	"perm": 227,
	"phag": 331,
	"phli": 131,
	"phlp": 132,
	"phlv": 133,
	"phnx": 115,
	"piqd": 293,
	"plrd": 282,
	"prti": 130,
	"psin": 103,
	"ranj": 303,
	"rjng": 363,
	"rohg": 167,
	"roro": 620,
	"runr": 211,
	"samr": 123,
	"sara": 292,
	"sarb": 105,
	"saur": 344,
	"sgnw": 95,
	"shaw": 281,
	"shrd": 319,
	"shui": 530,
	"sidd": 302,
	"sidt": 180,
	"sind": 318,
	"sinh": 348,
	"sogd": 141,
	"sogo": 142,
	"sora": 398,
	"soyo": 329,
	"sund": 362,
	"sunu": 274,
	"sylo": 316,
	"syrc": 135,
	"syre": 138,
	"syrj": 137,
	"syrn": 136,
	"tagb": 373,
	"takr": 321,
	"tale": 353,
	"talu": 354,
	"taml": 346,
	"tang": 520,
	"tavt": 359,
	"tayo": 380,
	"telu": 340,
	"teng": 290,
	"tfng": 120,
	"tglg": 370,
	"thaa": 170,
	"thai": 352,
	"tibt": 330,
	"tirh": 326,
	"tnsa": 275,
	"todr": 229,
	"tols": 299,
	"toto": 294,
	"tutg": 341,
	"ugar": 40,
	"vaii": 470,
	"visp": 280,
	"vith": 228,
	"wara": 262,
	"wcho": 283,
	"wole": 480,
	"xpeo": 30,
	"xsux": 20,
	"yezi": 192,
	"yiii": 460,
	"zanb": 339,
	"zinh": 994,
	"zmth": 995,
	"zsye": 993,
	"zsym": 996,
	"zxxx": 997,
	"zyyy": 998,
	"zzzz": 999,
}

type Language struct {
	Language string
	Script   string
	Region   string
}

// NewLanguage parses the string and returns an instance of Language struct.
func NewLanguage(s string) Language {
	language, err := ParseLanguage(s)
	if err != nil {
		return Language{}
	}

	return language
}

// ParseLanguage parses the given string as a language and returns it as a Language.
// If the string cannot be parsed an appropriate error is returned.
func ParseLanguage(s string) (Language, error) {
	// RFC 4647, 2.1 Basic Language Range
	language := Language{}
	var consumed bool
	if language.Language, language.Script, language.Region, s, consumed = consumeLanguageTags(skipWhitespaces(s)); !consumed {
		return Language{}, ErrInvalidLanguage
	}

	// there must not be anything left after parsing the header
	if len(s) > 0 {
		return Language{}, ErrInvalidMediaType
	}

	return language, nil
}

func consumeTag(s string) (language, remaining string, consumed bool) {
	// RFC 5646, 2.1. Syntax
	for i := 0; i < len(s); i++ {
		if !isAlphaChar(s[i]) {
			if len(s) >= 2 {
				return strings.ToLower(s[:i]), s[i:], true
			} else {
				return "", s, false
			}
		}
	}

	return strings.ToLower(s), "", len(s) >= 2
}

func consumeLanguageTags(s string) (language, script, region, remaining string, consumed bool) {
	language, s, consumed = consumeTag(s)

	if !consumed {
		return "", "", "", "", false
	}

	if len(language) < 2 {
		return "", "", "", s, false
	} else if len(language) == 2 {
		if _, ok := languageSet1[language]; !ok {
			return "", "", "", s, false
		}
	} else if len(language) == 3 {
		if _, ok := languageSet2[language]; !ok {
			return "", "", "", s, false
		}
	}

	if len(s) == 0 {
		return language, "", "", "", true
	}

	if s[0] != '-' {
		return "", "", "", "", false
	}

	tag1, s, consumed := consumeTag(s[1:])

	if len(tag1) == 4 {
		if _, ok := scripts[tag1]; !ok {
			return "", "", "", s, false
		} else {
			script = tag1
		}
	} else if len(tag1) == 2 {
		region = tag1
	}

	if len(s) == 0 {
		return language, script, region, s, true
	}

	if s[0] != '-' {
		return "", "", "", "", false
	}

	if len(region) == 0 {
		region, s, consumed = consumeTag(s[1:])

		return language, script, region, s, true
	}

	return language, script, region, s, true
}
