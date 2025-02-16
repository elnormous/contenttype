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
var scripts = map[string]string{
	"adlm": "166",
	"afak": "439",
	"aghb": "239",
	"ahom": "338",
	"arab": "160",
	"aran": "161",
	"armi": "124",
	"armn": "230",
	"avst": "134",
	"bali": "360",
	"bamu": "435",
	"bass": "259",
	"batk": "365",
	"beng": "325",
	"berf": "258",
	"bhks": "334",
	"blis": "550",
	"bopo": "285",
	"brah": "300",
	"brai": "570",
	"bugi": "367",
	"buhd": "372",
	"cakm": "349",
	"cans": "440",
	"cari": "201",
	"cham": "358",
	"cher": "445",
	"chis": "298",
	"chrs": "109",
	"cirt": "291",
	"copt": "204",
	"cpmn": "402",
	"cprt": "403",
	"cyrl": "220",
	"cyrs": "221",
	"deva": "315",
	"diak": "342",
	"dogr": "328",
	"dsrt": "250",
	"dupl": "755",
	"egyd": "070",
	"egyh": "060",
	"egyp": "050",
	"elba": "226",
	"elym": "128",
	"ethi": "430",
	"gara": "164",
	"geok": "241",
	"geor": "240",
	"glag": "225",
	"gong": "312",
	"gonm": "313",
	"goth": "206",
	"gran": "343",
	"grek": "200",
	"gujr": "320",
	"gukh": "397",
	"guru": "310",
	"hanb": "503",
	"hang": "286",
	"hani": "500",
	"hano": "371",
	"hans": "501",
	"hant": "502",
	"hatr": "127",
	"hebr": "125",
	"hira": "410",
	"hluw": "080",
	"hmng": "450",
	"hmnp": "451",
	"hrkt": "412",
	"hung": "176",
	"inds": "610",
	"ital": "210",
	"jamo": "284",
	"java": "361",
	"jpan": "413",
	"jurc": "510",
	"kali": "357",
	"kana": "411",
	"kawi": "368",
	"khar": "305",
	"khmr": "355",
	"khoj": "322",
	"kitl": "505",
	"kits": "288",
	"knda": "345",
	"kore": "287",
	"kpel": "436",
	"krai": "396",
	"kthi": "317",
	"lana": "351",
	"laoo": "356",
	"latf": "217",
	"latg": "216",
	"latn": "215",
	"leke": "364",
	"lepc": "335",
	"limb": "336",
	"lina": "400",
	"linb": "401",
	"lisu": "399",
	"loma": "437",
	"lyci": "202",
	"lydi": "116",
	"mahj": "314",
	"maka": "366",
	"mand": "140",
	"mani": "139",
	"marc": "332",
	"maya": "90",
	"medf": "265",
	"mend": "438",
	"merc": "101",
	"mero": "100",
	"mlym": "347",
	"modi": "324",
	"mong": "145",
	"moon": "218",
	"mroo": "264",
	"mtei": "337",
	"mult": "323",
	"mymr": "350",
	"nagm": "295",
	"nand": "311",
	"narb": "106",
	"nbat": "159",
	"newa": "333",
	"nkdb": "085",
	"nkgb": "420",
	"nkoo": "165",
	"nshu": "499",
	"ogam": "212",
	"olck": "261",
	"onao": "296",
	"orkh": "175",
	"orya": "327",
	"osge": "219",
	"osma": "260",
	"ougr": "143",
	"palm": "126",
	"pauc": "263",
	"pcun": "015",
	"pelm": "016",
	"perm": "227",
	"phag": "331",
	"phli": "131",
	"phlp": "132",
	"phlv": "133",
	"phnx": "115",
	"piqd": "293",
	"plrd": "282",
	"prti": "130",
	"psin": "103",
	"ranj": "303",
	"rjng": "363",
	"rohg": "167",
	"roro": "620",
	"runr": "211",
	"samr": "123",
	"sara": "292",
	"sarb": "105",
	"saur": "344",
	"sgnw": "095",
	"shaw": "281",
	"shrd": "319",
	"shui": "530",
	"sidd": "302",
	"sidt": "180",
	"sind": "318",
	"sinh": "348",
	"sogd": "141",
	"sogo": "142",
	"sora": "398",
	"soyo": "329",
	"sund": "362",
	"sunu": "274",
	"sylo": "316",
	"syrc": "135",
	"syre": "138",
	"syrj": "137",
	"syrn": "136",
	"tagb": "373",
	"takr": "321",
	"tale": "353",
	"talu": "354",
	"taml": "346",
	"tang": "520",
	"tavt": "359",
	"tayo": "380",
	"telu": "340",
	"teng": "290",
	"tfng": "120",
	"tglg": "370",
	"thaa": "170",
	"thai": "352",
	"tibt": "330",
	"tirh": "326",
	"tnsa": "275",
	"todr": "229",
	"tols": "299",
	"toto": "294",
	"tutg": "341",
	"ugar": "040",
	"vaii": "470",
	"visp": "280",
	"vith": "228",
	"wara": "262",
	"wcho": "283",
	"wole": "480",
	"xpeo": "030",
	"xsux": "020",
	"yezi": "192",
	"yiii": "460",
	"zanb": "339",
	"zinh": "994",
	"zmth": "995",
	"zsye": "993",
	"zsym": "996",
	"zxxx": "997",
	"zyyy": "998",
	"zzzz": "999",
}

// List of ISO 3166-1 countries
var countryCodes = map[string]string{
	"ad": "020",
	"ae": "784",
	"af": "004",
	"ag": "028",
	"ai": "660",
	"al": "008",
	"am": "051",
	"ao": "024",
	"aq": "010",
	"ar": "032",
	"as": "016",
	"at": "040",
	"au": "036",
	"aw": "533",
	"ax": "248",
	"az": "031",
	"ba": "070",
	"bb": "052",
	"bd": "050",
	"be": "056",
	"bf": "854",
	"bg": "100",
	"bh": "048",
	"bi": "108",
	"bj": "204",
	"bl": "652",
	"bm": "060",
	"bn": "096",
	"bo": "068",
	"bq": "535",
	"br": "076",
	"bs": "044",
	"bt": "064",
	"bv": "074",
	"bw": "072",
	"by": "112",
	"bz": "084",
	"ca": "124",
	"cc": "166",
	"cd": "180",
	"cf": "140",
	"cg": "178",
	"ch": "756",
	"ci": "384",
	"ck": "184",
	"cl": "152",
	"cm": "120",
	"cn": "156",
	"co": "170",
	"cr": "188",
	"cu": "192",
	"cv": "132",
	"cw": "531",
	"cx": "162",
	"cy": "196",
	"cz": "203",
	"de": "276",
	"dj": "262",
	"dk": "208",
	"dm": "212",
	"do": "214",
	"dz": "012",
	"ec": "218",
	"ee": "233",
	"eg": "818",
	"eh": "732",
	"er": "232",
	"es": "724",
	"et": "231",
	"fi": "246",
	"fj": "242",
	"fk": "238",
	"fm": "583",
	"fo": "234",
	"fr": "250",
	"ga": "266",
	"gb": "826",
	"gd": "308",
	"ge": "268",
	"gf": "254",
	"gg": "831",
	"gh": "288",
	"gi": "292",
	"gl": "304",
	"gm": "270",
	"gn": "324",
	"gp": "312",
	"gq": "226",
	"gr": "300",
	"gs": "239",
	"gt": "320",
	"gu": "316",
	"gw": "624",
	"gy": "328",
	"hk": "344",
	"hm": "334",
	"hn": "340",
	"hr": "191",
	"ht": "332",
	"hu": "348",
	"id": "360",
	"ie": "372",
	"il": "376",
	"im": "833",
	"in": "356",
	"io": "086",
	"iq": "368",
	"ir": "364",
	"is": "352",
	"it": "380",
	"je": "832",
	"jm": "388",
	"jo": "400",
	"jp": "392",
	"ke": "404",
	"kg": "417",
	"kh": "116",
	"ki": "296",
	"km": "174",
	"kn": "659",
	"kp": "408",
	"kr": "410",
	"kw": "414",
	"ky": "136",
	"kz": "398",
	"la": "418",
	"lb": "422",
	"lc": "662",
	"li": "438",
	"lk": "144",
	"lr": "430",
	"ls": "426",
	"lt": "440",
	"lu": "442",
	"lv": "428",
	"ly": "434",
	"ma": "504",
	"mc": "492",
	"md": "498",
	"me": "499",
	"mf": "663",
	"mg": "450",
	"mh": "584",
	"mk": "807",
	"ml": "466",
	"mm": "104",
	"mn": "496",
	"mo": "446",
	"mp": "580",
	"mq": "474",
	"mr": "478",
	"ms": "500",
	"mt": "470",
	"mu": "480",
	"mv": "462",
	"mw": "454",
	"mx": "484",
	"my": "458",
	"mz": "508",
	"na": "516",
	"nc": "540",
	"ne": "562",
	"nf": "574",
	"ng": "566",
	"ni": "558",
	"nl": "528",
	"no": "578",
	"np": "524",
	"nr": "520",
	"nu": "570",
	"nz": "554",
	"om": "512",
	"pa": "591",
	"pe": "604",
	"pf": "258",
	"pg": "598",
	"ph": "608",
	"pk": "586",
	"pl": "616",
	"pm": "666",
	"pn": "612",
	"pr": "630",
	"ps": "275",
	"pt": "620",
	"pw": "585",
	"py": "600",
	"qa": "634",
	"re": "638",
	"ro": "642",
	"rs": "688",
	"ru": "643",
	"rw": "646",
	"sa": "682",
	"sb": "090",
	"sc": "690",
	"sd": "729",
	"se": "752",
	"sg": "702",
	"sh": "654",
	"si": "705",
	"sj": "744",
	"sk": "703",
	"sl": "694",
	"sm": "674",
	"sn": "686",
	"so": "706",
	"sr": "740",
	"ss": "728",
	"st": "678",
	"sv": "222",
	"sx": "534",
	"sy": "760",
	"sz": "748",
	"tc": "796",
	"td": "148",
	"tf": "260",
	"tg": "768",
	"th": "764",
	"tj": "762",
	"tk": "772",
	"tl": "626",
	"tm": "795",
	"tn": "788",
	"to": "776",
	"tr": "792",
	"tt": "780",
	"tv": "798",
	"tw": "158",
	"tz": "834",
	"ua": "804",
	"ug": "800",
	"um": "581",
	"us": "840",
	"uy": "858",
	"uz": "860",
	"va": "336",
	"vc": "670",
	"ve": "862",
	"vg": "092",
	"vi": "850",
	"vn": "704",
	"vu": "548",
	"wf": "876",
	"ws": "882",
	"ye": "887",
	"yt": "175",
	"za": "710",
	"zm": "894",
	"zw": "716",
}

// List of ISO 3166-1 countries
var countryNumbers = map[string]string{
	"004": "af",
	"008": "al",
	"010": "aq",
	"012": "dz",
	"016": "as",
	"020": "ad",
	"024": "ao",
	"028": "ag",
	"031": "az",
	"032": "ar",
	"036": "au",
	"040": "at",
	"044": "bs",
	"048": "bh",
	"050": "bd",
	"051": "am",
	"052": "bb",
	"056": "be",
	"060": "bm",
	"064": "bt",
	"068": "bo",
	"070": "ba",
	"072": "bw",
	"074": "bv",
	"076": "br",
	"084": "bz",
	"086": "io",
	"090": "sb",
	"092": "vg",
	"096": "bn",
	"100": "bg",
	"104": "mm",
	"108": "bi",
	"112": "by",
	"116": "kh",
	"120": "cm",
	"124": "ca",
	"132": "cv",
	"136": "ky",
	"140": "cf",
	"144": "lk",
	"148": "td",
	"152": "cl",
	"156": "cn",
	"158": "tw",
	"162": "cx",
	"166": "cc",
	"170": "co",
	"174": "km",
	"175": "yt",
	"178": "cg",
	"180": "cd",
	"184": "ck",
	"188": "cr",
	"191": "hr",
	"192": "cu",
	"196": "cy",
	"203": "cz",
	"204": "bj",
	"208": "dk",
	"212": "dm",
	"214": "do",
	"218": "ec",
	"222": "sv",
	"226": "gq",
	"231": "et",
	"232": "er",
	"233": "ee",
	"234": "fo",
	"238": "fk",
	"239": "gs",
	"242": "fj",
	"246": "fi",
	"248": "ax",
	"250": "fr",
	"254": "gf",
	"258": "pf",
	"260": "tf",
	"262": "dj",
	"266": "ga",
	"268": "ge",
	"270": "gm",
	"275": "ps",
	"276": "de",
	"288": "gh",
	"292": "gi",
	"296": "ki",
	"300": "gr",
	"304": "gl",
	"308": "gd",
	"312": "gp",
	"316": "gu",
	"320": "gt",
	"324": "gn",
	"328": "gy",
	"332": "ht",
	"334": "hm",
	"336": "va",
	"340": "hn",
	"344": "hk",
	"348": "hu",
	"352": "is",
	"356": "in",
	"360": "id",
	"364": "ir",
	"368": "iq",
	"372": "ie",
	"376": "il",
	"380": "it",
	"384": "ci",
	"388": "jm",
	"392": "jp",
	"398": "kz",
	"400": "jo",
	"404": "ke",
	"408": "kp",
	"410": "kr",
	"414": "kw",
	"417": "kg",
	"418": "la",
	"422": "lb",
	"426": "ls",
	"428": "lv",
	"430": "lr",
	"434": "ly",
	"438": "li",
	"440": "lt",
	"442": "lu",
	"446": "mo",
	"450": "mg",
	"454": "mw",
	"458": "my",
	"462": "mv",
	"466": "ml",
	"470": "mt",
	"474": "mq",
	"478": "mr",
	"480": "mu",
	"484": "mx",
	"492": "mc",
	"496": "mn",
	"498": "md",
	"499": "me",
	"500": "ms",
	"504": "ma",
	"508": "mz",
	"512": "om",
	"516": "na",
	"520": "nr",
	"524": "np",
	"528": "nl",
	"531": "cw",
	"533": "aw",
	"534": "sx",
	"535": "bq",
	"540": "nc",
	"548": "vu",
	"554": "nz",
	"558": "ni",
	"562": "ne",
	"566": "ng",
	"570": "nu",
	"574": "nf",
	"578": "no",
	"580": "mp",
	"581": "um",
	"583": "fm",
	"584": "mh",
	"585": "pw",
	"586": "pk",
	"591": "pa",
	"598": "pg",
	"600": "py",
	"604": "pe",
	"608": "ph",
	"612": "pn",
	"616": "pl",
	"620": "pt",
	"624": "gw",
	"626": "tl",
	"630": "pr",
	"634": "qa",
	"638": "re",
	"642": "ro",
	"643": "ru",
	"646": "rw",
	"652": "bl",
	"654": "sh",
	"659": "kn",
	"660": "ai",
	"662": "lc",
	"663": "mf",
	"666": "pm",
	"670": "vc",
	"674": "sm",
	"678": "st",
	"682": "sa",
	"686": "sn",
	"688": "rs",
	"690": "sc",
	"694": "sl",
	"702": "sg",
	"703": "sk",
	"704": "vn",
	"705": "si",
	"706": "so",
	"710": "za",
	"716": "zw",
	"724": "es",
	"728": "ss",
	"729": "sd",
	"732": "eh",
	"740": "sr",
	"744": "sj",
	"748": "sz",
	"752": "se",
	"756": "ch",
	"760": "sy",
	"762": "tj",
	"764": "th",
	"768": "tg",
	"772": "tk",
	"776": "to",
	"780": "tt",
	"784": "ae",
	"788": "tn",
	"792": "tr",
	"795": "tm",
	"796": "tc",
	"798": "tv",
	"800": "ug",
	"804": "ua",
	"807": "mk",
	"818": "eg",
	"826": "gb",
	"831": "gg",
	"832": "je",
	"833": "im",
	"834": "tz",
	"840": "us",
	"850": "vi",
	"854": "bf",
	"858": "uy",
	"860": "uz",
	"862": "ve",
	"876": "wf",
	"882": "ws",
	"887": "ye",
	"894": "zm",
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

func consumeLanguageTag(s string) (language, remaining string, consumed bool) {
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

func consumeTag(s string) (language, remaining string, consumed bool) {
	// RFC 5646, 2.1. Syntax
	for i := 0; i < len(s); i++ {
		if !isAlphaChar(s[i]) && !isDigitChar(s[i]) {
			if len(s) >= 2 {
				return strings.ToLower(s[:i]), s[i:], true
			} else {
				return "", s, false
			}
		}
	}

	return strings.ToLower(s), "", len(s) >= 2
}

func isValidLanguage(language string) bool {
	if len(language) == 2 {
		_, ok := languageSet1[language]
		return ok
	} else if len(language) == 3 {
		_, ok := languageSet2[language]
		return ok
	}

	return false
}

func isValidScript(script string) bool {
	_, ok := scripts[script]
	return ok
}

func isValidCountry(country string) bool {
	if len(country) == 2 {
		_, ok := countryCodes[country]
		return ok
	} else if len(country) == 3 {
		_, ok := countryNumbers[country]
		return ok
	}

	return false
}

func consumeLanguageTags(s string) (language, script, region, remaining string, consumed bool) {
	language, s, consumed = consumeLanguageTag(s)

	if !consumed {
		return "", "", "", "", false
	}

	if !isValidLanguage(language) {
		return "", "", "", s, false
	}

	if len(s) == 0 {
		return language, "", "", "", true
	}

	if s[0] != '-' {
		return "", "", "", "", false
	}

	s = s[1:]

	tag1, s, consumed := consumeTag(s)

	if !consumed {
		return "", "", "", "", false
	}

	if len(tag1) == 4 {
		if !isValidScript(tag1) {
			return "", "", "", s, false
		} else {
			script = tag1
		}
	} else if len(tag1) == 3 || len(tag1) == 2 {
		if !isValidCountry(tag1) {
			return "", "", "", s, false
		} else {
			region = tag1
		}
	} else {
		return "", "", "", s, false
	}

	if len(s) == 0 {
		return language, script, region, s, true
	}

	if s[0] != '-' {
		return "", "", "", "", false
	}

	s = s[1:]

	if len(region) == 0 {
		region, s, consumed = consumeTag(s)

		if len(region) == 3 || len(region) == 2 {
			if !isValidCountry(region) {
				return "", "", "", s, false
			}
		} else {
			return "", "", "", s, false
		}
	}

	if len(s) == 0 {
		return language, script, region, s, true
	}

	return language, script, region, s, true
}
