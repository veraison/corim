// Copyright 2021 Contributors to the Veraison project.
// SPDX-License-Identifier: Apache-2.0

package cmd

import "github.com/veraison/corim/comid"

var (
	minimalCorimTemplate = []byte(`{
		"corim-id": "5c57e8f4-46cd-421b-91c9-08cf93e13cfc"
	}`)
	badCBOR = comid.MustHexDecode(nil, "ffff")
	// a "tag-id only" CoMID {1: {0: h'366D0A0A598845ED84882F2A544F6242'}}
	invalidComid = comid.MustHexDecode(nil,
		"a101a10050366d0a0a598845ed84882f2a544f6242",
	)
	//
	invalidCots = comid.MustHexDecode(nil, "a2028006a100f6")
	// note: embedded CoSWIDs are not validated {0: h'5C57E8F446CD421B91C908CF93E13CFC', 1: [505(h'deadbeef')]}
	testCorimValid = comid.MustHexDecode(nil,
		"a200505c57e8f446cd421b91c908cf93e13cfc0181d901f944deadbeef",
	)
	// {0: h'5C57E8F446CD421B91C908CF93E13CFC'}
	testCorimInvalid = comid.MustHexDecode(nil,
		"a100505c57e8f446cd421b91c908cf93e13cfc",
	)
	testMetaInvalid = []byte("{}")
	testMetaValid   = []byte(`{
		"signer": {
			"name": "ACME Ltd signing key",
			"uri": "https://acme.example"
		},
		"validity": {
			"not-before": "2021-12-31T00:00:00Z",
			"not-after": "2025-12-31T00:00:00Z"
		}
	}`)
	testECKey = []byte(`{
		"kty": "EC",
		"crv": "P-256",
		"x": "MKBCTNIcKUSDii11ySs3526iDZ8AiTo7Tu6KPAqv7D4",
		"y": "4Etl6SRW2YiLUrN5vfvVHuhp7x8PxltmWWlbbM4IFyM",
		"d": "870MB6gfuTJ4HtUnUvYMyJpr5eUZNP4Bk43bVdj3eAE",
		"use": "enc",
		"kid": "1"
	  }`)

	testSignedCorimValid = comid.MustHexDecode(nil, `
	d284585da3012603746170706c69636174696f6e2f72696d2b63626f7208
	5841a200a2007441434d45204c7464207369676e696e67206b657901d820
	7468747470733a2f2f61636d652e6578616d706c6501a200c11a61ce4800
	01c11a69546780a0590282a600505c57e8f446cd421b91c908cf93e13cfc
	01815901a3d901faa40065656e2d474201a1005043bbe37f2e614b33aed3
	53cff1428b160281a3006941434d45204c74642e01d8207468747470733a
	2f2f61636d652e6578616d706c65028300010204a1008182a100a300d902
	58582061636d652d696d706c656d656e746174696f6e2d69642d30303030
	3030303031016441434d45026a526f616452756e6e657283a200d90259a3
	0162424c0465322e312e30055820acbb11c7e4da217205523ce4ce1a245a
	e1a239ae3c6bfd9e7871f7e5d8bae86b01a102818201582087428fc52280
	3d31065e7bce3cf03fe475096631e5e07bbd7a0fde60c4cf25c7a200d902
	59a3016450526f540465312e332e35055820acbb11c7e4da217205523ce4
	ce1a245ae1a239ae3c6bfd9e7871f7e5d8bae86b01a10281820158200263
	829989b6fd954f72baaf2fc64bc2e2f01d692d4de72986ea808f6e99813f
	a200d90259a3016441526f540465302e312e34055820acbb11c7e4da2172
	05523ce4ce1a245ae1a239ae3c6bfd9e7871f7e5d8bae86b01a102818201
	5820a3a5e715f0cc574a73c3f9bebb6bc24f32ffd5b67b387244c2c909da
	779a14780281a200d820784068747470733a2f2f706172656e742e657861
	6d706c652f72696d732f63636233616138352d363162342d343066312d38
	3438652d3032616436653861323534620182015820e45b72f5c0c0b572db
	4d8d3ab7e97f368ff74e62347a824decb67a84e5224d750382482b060104
	01a02064781c687474703a2f2f61726d2e636f6d2f696f742f70726f6669
	6c652f3104a200c11a61ce480001c11a695467800581a3006941434d4520
	4c74642e01d8206c61636d652e6578616d706c6502810158400f79b7a7de
	b0480ed9c5b8b29aaffef33b1094f4e11409dbbbc4c1c2abf79fda30d72e
	b86d19053dc5caada1abe76f1e13cad3bef27eebc09086463fe02f5ac4	
	`)
	testSignedCorimInvalid = comid.MustHexDecode(nil, `
d284585da3012603746170706c69636174696f6e2f72696d2b63626f7208
5841a200a2007441434d45204c7464207369676e696e67206b657901d820
7468747470733a2f2f61636d652e6578616d706c6501a200c11a61ce4800
01c11a69546780a041a044deadbeef
	`)

	testSignedCorimValidWithCots = comid.MustHexDecode(nil, `
	d2845835a3012603746170706c69636174696f6e2f72696d2b63626f
	72085819a100a1007441434d45204c7464207369676e696e67206b65
	79a05908aea200505c57e8f446cd421b91c908cf93e13cfc01815908
	96d901fba20281a101a100a100d86f442a03040506a1008482015902
	d9a28202d5308202d13059301306072a8648ce3d020106082a8648ce
	3d03010703420004cdd1fe64cf2c04cd93986da576dcedcdf1c92edf
	d2682cd8e51cfc0409b5c30d6a742e900ed73db8f3f868cba516fd36
	4c4cf3d1ffddd7c07b06d7a992172f1704148a84cff98095a3bc36d6
	eea518d6978d9bd71f603082025c305c310b300906035504060c0255
	53311f301d060355040a0c16536e6f6262697368204170706172656c
	2c20496e632e312c302a06035504030c23536e6f6262697368204170
	706172656c2c20496e632e20547275737420416e63686f72a08201fa
	3082019fa0030201020214101b934465c01045441e1bb8c5a7c09ea9
	bea988300a06082a8648ce3d040302305c310b300906035504060c02
	5553311f301d060355040a0c16536e6f626269736820417070617265
	6c2c20496e632e312c302a06035504030c23536e6f62626973682041
	70706172656c2c20496e632e20547275737420416e63686f72301e17
	0d3232303531393135313330385a170d333230353136313531333038
	5a305c310b300906035504060c025553311f301d060355040a0c1653
	6e6f6262697368204170706172656c2c20496e632e312c302a060355
	04030c23536e6f6262697368204170706172656c2c20496e632e2054
	7275737420416e63686f723059301306072a8648ce3d020106082a86
	48ce3d03010703420004cdd1fe64cf2c04cd93986da576dcedcdf1c9
	2edfd2682cd8e51cfc0409b5c30d6a742e900ed73db8f3f868cba516
	fd364c4cf3d1ffddd7c07b06d7a992172f17a33f303d301d0603551d
	0e041604148a84cff98095a3bc36d6eea518d6978d9bd71f60300b06
	03551d0f040403020284300f0603551d130101ff040530030101ff30
	0a06082a8648ce3d0403020349003046022100b671fe377f73cf9423
	bafddc6fe347ed220c714e828717bd94cc436fefecb8ca02210081ad
	0bfafa48668531a3fbcc3d8496806e2174edc96bbc1f097e60675745
	8f7782015902baa28202b6308202b23059301306072a8648ce3d0201
	06082a8648ce3d0301070342000497cf6d70d76a30400c79f1ebab6a
	d6168871248710d4f4c1207adaef02d19867136791f2c2b29f6f2dcc
	3c2125c2f2d2fdee4f59094b67af43332f0dfb5ca4400414f6dad1e5
	128bbf0de9e95343b371c6f7ffe7e26e3082023d3052310b30090603
	5504060c025553311a3018060355040a0c115a657374792048616e64
	732c20496e632e3127302506035504030c1e5a657374792048616e64
	732c20496e632e20547275737420416e63686f72a08201e53082018b
	a00302010202140bdc4aa05179503e58f275d552467347bcafb53330
	0a06082a8648ce3d0403023052310b300906035504060c025553311a
	3018060355040a0c115a657374792048616e64732c20496e632e3127
	302506035504030c1e5a657374792048616e64732c20496e632e2054
	7275737420416e63686f72301e170d3232303531393135313330375a
	170d3332303531363135313330375a3052310b300906035504060c02
	5553311a3018060355040a0c115a657374792048616e64732c20496e
	632e3127302506035504030c1e5a657374792048616e64732c20496e
	632e20547275737420416e63686f723059301306072a8648ce3d0201
	06082a8648ce3d0301070342000497cf6d70d76a30400c79f1ebab6a
	d6168871248710d4f4c1207adaef02d19867136791f2c2b29f6f2dcc
	3c2125c2f2d2fdee4f59094b67af43332f0dfb5ca440a33f303d301d
	0603551d0e04160414f6dad1e5128bbf0de9e95343b371c6f7ffe7e2
	6e300b0603551d0f040403020284300f0603551d130101ff04053003
	0101ff300a06082a8648ce3d040302034800304502201da58be7fa44
	2c6cd849ef3567224a9203c225158966b21afd40f3192cf347980221
	009827b3e0a1aba2514a3994fa6efa9fd6c610b8905fbed93fcb5250
	754be8aa58820159027ea282027a308202763059301306072a8648ce
	3d020106082a8648ce3d03010703420004e351aa10392407a4bd037c
	a0bc1154d0e701f0674f39cb2c4f9210692cebbeec1d27977dc56165
	751e0e237bfdfb1536e99a884591429661df35cfc0bf2b50cc041401
	5c45c9acb0462a715dd710a078c01549f1013f30820201303e310b30
	0906035504060c0255533110300e060355040a0c074578616d706c65
	311d301b06035504030c144578616d706c6520547275737420416e63
	686f72a08201bd30820164a003020102021500d09d90bf3d525cc773
	d522ed77d59e22bba45b88300a06082a8648ce3d040302303e310b30
	0906035504060c0255533110300e060355040a0c074578616d706c65
	311d301b06035504030c144578616d706c6520547275737420416e63
	686f72301e170d3232303531393135313330375a170d333230353136
	3135313330375a303e310b300906035504060c0255533110300e0603
	55040a0c074578616d706c65311d301b06035504030c144578616d70
	6c6520547275737420416e63686f723059301306072a8648ce3d0201
	06082a8648ce3d03010703420004e351aa10392407a4bd037ca0bc11
	54d0e701f0674f39cb2c4f9210692cebbeec1d27977dc56165751e0e
	237bfdfb1536e99a884591429661df35cfc0bf2b50cca33f303d301d
	0603551d0e04160414015c45c9acb0462a715dd710a078c01549f101
	3f300b0603551d0f040403020284300f0603551d130101ff04053003
	0101ff300a06082a8648ce3d040302034700304402200b06274006c7
	e9cc6d254cbf4487d6fa0146ea9f317e9281196e0be9a6fbedf50220
	56813e6e110dd23eb048fcde3e32eb11d0fe3c48328c7279adb035d5
	23eaff538202585b3059301306072a8648ce3d020106082a8648ce3d
	03010703420004ad8a0c01da9eda0253dc2bc27227d9c7213df8df13
	e89cb9cdb7a8e4b62d9ce8a99a2d705c0f7f80db65c006d1091422b4
	7fc611cbd46869733d9c483884d5fe5840fd190f51a1f11dd477335e
	82b0304e9f5987d43c921cb6ec43e23dda1603ff6f1f3bdab11a5467
	039d0ac0901cf2d344001c4fcf12629ea116ebf758aae0ec79
	`)
	PSARefValCBOR = comid.MustHexDecode(nil, `
a40065656e2d474201a1005043bbe37f2e614b33aed353cff1428b160281
a3006941434d45204c74642e01d8207468747470733a2f2f61636d652e65
78616d706c65028300010204a1008182a100a300d90258582061636d652d
696d706c656d656e746174696f6e2d69642d303030303030303031016441
434d45026a526f616452756e6e657283a200d90259a30162424c0465322e
312e30055820acbb11c7e4da217205523ce4ce1a245ae1a239ae3c6bfd9e
7871f7e5d8bae86b01a102818201582087428fc522803d31065e7bce3cf0
3fe475096631e5e07bbd7a0fde60c4cf25c7a200d90259a3016450526f54
0465312e332e35055820acbb11c7e4da217205523ce4ce1a245ae1a239ae
3c6bfd9e7871f7e5d8bae86b01a10281820158200263829989b6fd954f72
baaf2fc64bc2e2f01d692d4de72986ea808f6e99813fa200d90259a30164
41526f540465302e312e34055820acbb11c7e4da217205523ce4ce1a245a
e1a239ae3c6bfd9e7871f7e5d8bae86b01a1028182015820a3a5e715f0cc
574a73c3f9bebb6bc24f32ffd5b67b387244c2c909da779a1478
	`)
	testComid = comid.MustHexDecode(nil, `
a40065656e2d474201a1005043bbe37f2e614b33aed353cff1428b160281
a3006941434d45204c74642e01d8207468747470733a2f2f61636d652e65
78616d706c65028300010204a1008182a100a300d90258582061636d652d
696d706c656d656e746174696f6e2d69642d303030303030303031016441
434d45026a526f616452756e6e657283a200d90259a30162424c0465322e
312e30055820acbb11c7e4da217205523ce4ce1a245ae1a239ae3c6bfd9e
7871f7e5d8bae86b01a102818201582087428fc522803d31065e7bce3cf0
3fe475096631e5e07bbd7a0fde60c4cf25c7a200d90259a3016450526f54
0465312e332e35055820acbb11c7e4da217205523ce4ce1a245ae1a239ae
3c6bfd9e7871f7e5d8bae86b01a10281820158200263829989b6fd954f72
baaf2fc64bc2e2f01d692d4de72986ea808f6e99813fa200d90259a30164
41526f540465302e312e34055820acbb11c7e4da217205523ce4ce1a245a
e1a239ae3c6bfd9e7871f7e5d8bae86b01a1028182015820a3a5e715f0cc
574a73c3f9bebb6bc24f32ffd5b67b387244c2c909da779a1478
	`)
	testCoswid = comid.MustHexDecode(nil, `
a8007820636f6d2e61636d652e727264323031332d63652d7370312d7634
2d312d352d300c0001783041434d4520526f616472756e6e657220446574
6563746f72203230313320436f796f74652045646974696f6e205350310d
65342e312e3505a5182b65747269616c182d6432303133182f66636f796f
7465183473526f616472756e6e6572204465746563746f72183663737031
0282a3181f745468652041434d4520436f72706f726174696f6e18206861
636d652e636f6d1821820102a3181f75436f796f74652053657276696365
732c20496e632e18206c6d79636f796f74652e636f6d18210404a2182678
1c7777772e676e752e6f72672f6c6963656e7365732f67706c2e74787418
28676c6963656e736506a110a318186a72726465746563746f7218196d25
70726f6772616d6461746125181aa111a318186e72726465746563746f72
2e657865141a000820e80782015820a314fc2dc663ae7a6b6bc678759405
7396e6b3f569cd50fd5ddb4d1bbafd2b6a
	`)

	testCots = comid.MustHexDecode(nil, `
	a301a10050ab0f44b1bfdc4604ab4a30f80407ebcc0281a103764d697
	363656c6c616e656f75732054412053746f726506a1008382005901c1
	308201bd30820164a003020102021500d09d90bf3d525cc773d522ed7
	7d59e22bba45b88300a06082a8648ce3d040302303e310b3009060355
	04060c0255533110300e060355040a0c074578616d706c65311d301b0
	6035504030c144578616d706c6520547275737420416e63686f72301e
	170d3232303531393135313330375a170d33323035313631353133303
	75a303e310b300906035504060c0255533110300e060355040a0c0745
	78616d706c65311d301b06035504030c144578616d706c65205472757
	37420416e63686f723059301306072a8648ce3d020106082a8648ce3d
	03010703420004e351aa10392407a4bd037ca0bc1154d0e701f0674f3
	9cb2c4f9210692cebbeec1d27977dc56165751e0e237bfdfb1536e99a
	884591429661df35cfc0bf2b50cca33f303d301d0603551d0e0416041
	4015c45c9acb0462a715dd710a078c01549f1013f300b0603551d0f04
	0403020284300f0603551d130101ff040530030101ff300a06082a864
	8ce3d040302034700304402200b06274006c7e9cc6d254cbf4487d6fa
	0146ea9f317e9281196e0be9a6fbedf5022056813e6e110dd23eb048f
	cde3e32eb11d0fe3c48328c7279adb035d523eaff5382015902baa282
	02b6308202b23059301306072a8648ce3d020106082a8648ce3d03010
	70342000497cf6d70d76a30400c79f1ebab6ad6168871248710d4f4c1
	207adaef02d19867136791f2c2b29f6f2dcc3c2125c2f2d2fdee4f590
	94b67af43332f0dfb5ca4400414f6dad1e5128bbf0de9e95343b371c6
	f7ffe7e26e3082023d3052310b300906035504060c025553311a30180
	60355040a0c115a657374792048616e64732c20496e632e3127302506
	035504030c1e5a657374792048616e64732c20496e632e20547275737
	420416e63686f72a08201e53082018ba00302010202140bdc4aa05179
	503e58f275d552467347bcafb533300a06082a8648ce3d04030230523
	10b300906035504060c025553311a3018060355040a0c115a65737479
	2048616e64732c20496e632e3127302506035504030c1e5a657374792
	048616e64732c20496e632e20547275737420416e63686f72301e170d
	3232303531393135313330375a170d3332303531363135313330375a3
	052310b300906035504060c025553311a3018060355040a0c115a6573
	74792048616e64732c20496e632e3127302506035504030c1e5a65737
	4792048616e64732c20496e632e20547275737420416e63686f723059
	301306072a8648ce3d020106082a8648ce3d0301070342000497cf6d7
	0d76a30400c79f1ebab6ad6168871248710d4f4c1207adaef02d19867
	136791f2c2b29f6f2dcc3c2125c2f2d2fdee4f59094b67af43332f0df
	b5ca440a33f303d301d0603551d0e04160414f6dad1e5128bbf0de9e9
	5343b371c6f7ffe7e26e300b0603551d0f040403020284300f0603551
	d130101ff040530030101ff300a06082a8648ce3d0403020348003045
	02201da58be7fa442c6cd849ef3567224a9203c225158966b21afd40f
	3192cf347980221009827b3e0a1aba2514a3994fa6efa9fd6c610b890
	5fbed93fcb5250754be8aa5882015902d9a28202d5308202d13059301
	306072a8648ce3d020106082a8648ce3d03010703420004cdd1fe64cf
	2c04cd93986da576dcedcdf1c92edfd2682cd8e51cfc0409b5c30d6a7
	42e900ed73db8f3f868cba516fd364c4cf3d1ffddd7c07b06d7a99217
	2f1704148a84cff98095a3bc36d6eea518d6978d9bd71f603082025c3
	05c310b300906035504060c025553311f301d060355040a0c16536e6f
	6262697368204170706172656c2c20496e632e312c302a06035504030
	c23536e6f6262697368204170706172656c2c20496e632e2054727573
	7420416e63686f72a08201fa3082019fa0030201020214101b934465c
	01045441e1bb8c5a7c09ea9bea988300a06082a8648ce3d040302305c
	310b300906035504060c025553311f301d060355040a0c16536e6f626
	2697368204170706172656c2c20496e632e312c302a06035504030c23
	536e6f6262697368204170706172656c2c20496e632e2054727573742
	0416e63686f72301e170d3232303531393135313330385a170d333230
	3531363135313330385a305c310b300906035504060c025553311f301
	d060355040a0c16536e6f6262697368204170706172656c2c20496e63
	2e312c302a06035504030c23536e6f6262697368204170706172656c2
	c20496e632e20547275737420416e63686f723059301306072a8648ce
	3d020106082a8648ce3d03010703420004cdd1fe64cf2c04cd93986da
	576dcedcdf1c92edfd2682cd8e51cfc0409b5c30d6a742e900ed73db8
	f3f868cba516fd364c4cf3d1ffddd7c07b06d7a992172f17a33f303d3
	01d0603551d0e041604148a84cff98095a3bc36d6eea518d6978d9bd7
	1f60300b0603551d0f040403020284300f0603551d130101ff0405300
	30101ff300a06082a8648ce3d0403020349003046022100b671fe377f
	73cf9423bafddc6fe347ed220c714e828717bd94cc436fefecb8ca022
	10081ad0bfafa48668531a3fbcc3d8496806e2174edc96bbc1f097e60
	6757458f77
	`)
)
