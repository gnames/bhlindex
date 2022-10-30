# Stats for just walking BHL tree for different computers

## sf23 (IOPS: R:140k, W:59k)

```bash
10:00:34 INF Processing BHL items items=10,000 items/hour=669,803
10:01:30 INF Processing BHL items items=20,000 items/hour=660,566
10:02:26 INF Processing BHL items items=30,000 items/hour=653,866
10:03:20 INF Processing BHL items items=40,000 items/hour=656,198
10:04:08 INF Processing BHL items items=50,000 items/hour=672,934
10:04:58 INF Processing BHL items items=60,000 items/hour=680,288
10:05:52 INF Processing BHL items items=70,000 items/hour=678,264
10:06:57 INF Processing BHL items items=80,000 items/hour=660,324
10:07:39 INF Processing BHL items items=90,000 items/hour=677,596
10:08:31 INF Processing BHL items items=100,000 items/hour=678,768
10:09:21 INF Processing BHL items items=110,000 items/hour=681,913
10:09:46 INF Processing BHL items items=120,000 items/hour=714,069
10:10:11 INF Processing BHL items items=130,000 items/hour=742,724
10:10:25 INF Processing BHL items items=140,000 items/hour=782,001
10:10:38 INF Processing BHL items items=150,000 items/hour=821,612
10:10:51 INF Processing BHL items items=160,000 items/hour=859,412
10:11:01 INF Processing BHL items items=170,000 items/hour=898,874
10:11:14 INF Processing BHL items items=180,000 items/hour=934,915
10:11:29 INF Processing BHL items items=190,000 items/hour=965,002
10:11:52 INF Processing BHL items items=200,000 items/hour=984,649
10:12:21 INF Processing BHL items items=210,000 items/hour=994,642
10:12:55 INF Processing BHL items items=220,000 items/hour=996,958
10:13:38 INF Processing BHL items items=230,000 items/hour=989,357
```

## Dell XPS 13 2020 (IOPS: R:90, W:30)

```bash
04:24:14 INF Processing BHL items items=10,000 items/hour=520,715
04:25:18 INF Processing BHL items items=20,000 items/hour=539,564
04:26:22 INF Processing BHL items items=30,000 items/hour=546,763
04:27:24 INF Processing BHL items items=40,000 items/hour=555,056
04:28:22 INF Processing BHL items items=50,000 items/hour=567,609
04:29:23 INF Processing BHL items items=60,000 items/hour=570,744
04:30:27 INF Processing BHL items items=70,000 items/hour=569,912
04:31:37 INF Processing BHL items items=80,000 items/hour=562,168
04:32:25 INF Processing BHL items items=90,000 items/hour=579,001
04:33:25 INF Processing BHL items items=100,000 items/hour=580,934
04:34:22 INF Processing BHL items items=110,000 items/hour=584,839
04:34:47 INF Processing BHL items items=120,000 items/hour=614,952
04:35:15 INF Processing BHL items items=130,000 items/hour=641,426
04:35:30 INF Processing BHL items items=140,000 items/hour=676,801
04:35:43 INF Processing BHL items items=150,000 items/hour=712,560
04:35:56 INF Processing BHL items items=160,000 items/hour=747,333
04:36:06 INF Processing BHL items items=170,000 items/hour=783,178
04:36:19 INF Processing BHL items items=180,000 items/hour=816,296
04:36:36 INF Processing BHL items items=190,000 items/hour=843,843
04:36:59 INF Processing BHL items items=200,000 items/hour=863,225
04:37:30 INF Processing BHL items items=210,000 items/hour=873,980
04:38:08 INF Processing BHL items items=220,000 items/hour=876,983
04:38:55 INF Processing BHL items items=230,000 items/hour=871,212
```

## Gygabyte Aero 15 laptop NoCOW (IOPS: R:98k, W:32k)

```bash
06:21:13 INF Processing BHL items items=10,000 items/hour=504,582
06:22:23 INF Processing BHL items items=20,000 items/hour=507,947
06:23:33 INF Processing BHL items items=30,000 items/hour=511,682
06:24:45 INF Processing BHL items items=40,000 items/hour=508,712
06:25:50 INF Processing BHL items items=50,000 items/hour=516,593
06:26:57 INF Processing BHL items items=60,000 items/hour=520,336
06:28:04 INF Processing BHL items items=70,000 items/hour=522,427
06:29:22 INF Processing BHL items items=80,000 items/hour=514,368
06:30:14 INF Processing BHL items items=90,000 items/hour=528,716
06:31:22 INF Processing BHL items items=100,000 items/hour=528,786
06:32:26 INF Processing BHL items items=110,000 items/hour=531,966
06:32:55 INF Processing BHL items items=120,000 items/hour=558,607
06:33:26 INF Processing BHL items items=130,000 items/hour=582,139
06:33:43 INF Processing BHL items items=140,000 items/hour=613,872
06:33:58 INF Processing BHL items items=150,000 items/hour=646,029
06:34:12 INF Processing BHL items items=160,000 items/hour=677,051
06:34:24 INF Processing BHL items items=170,000 items/hour=709,383
06:34:38 INF Processing BHL items items=180,000 items/hour=739,083
06:34:58 INF Processing BHL items items=190,000 items/hour=762,774
06:35:26 INF Processing BHL items items=200,000 items/hour=779,253
06:36:00 INF Processing BHL items items=210,000 items/hour=788,569
06:36:43 INF Processing BHL items items=220,000 items/hour=791,234
06:37:37 INF Processing BHL items items=230,000 items/hour=784,645
```

## Gygabyte Aero 15 laptop COW+zstd:1 (IOPS: R:75k, W:25k)

```bash
04:44:12 INF Processing BHL items items=10,000 items/hour=340,792
04:45:58 INF Processing BHL items items=20,000 items/hour=339,417
04:47:45 INF Processing BHL items items=30,000 items/hour=339,370
04:49:32 INF Processing BHL items items=40,000 items/hour=338,363
04:51:05 INF Processing BHL items items=50,000 items/hour=347,301
04:52:39 INF Processing BHL items items=60,000 items/hour=352,348
04:54:20 INF Processing BHL items items=70,000 items/hour=353,120
04:56:11 INF Processing BHL items items=80,000 items/hour=349,271
04:57:27 INF Processing BHL items items=90,000 items/hour=359,859
04:59:05 INF Processing BHL items items=100,000 items/hour=360,487
05:00:35 INF Processing BHL items items=110,000 items/hour=363,759
05:01:17 INF Processing BHL items items=120,000 items/hour=381,999
05:02:04 INF Processing BHL items items=130,000 items/hour=397,426
05:02:31 INF Processing BHL items items=140,000 items/hour=418,319
05:02:56 INF Processing BHL items items=150,000 items/hour=439,169
05:03:19 INF Processing BHL items items=160,000 items/hour=459,769
05:03:38 INF Processing BHL items items=170,000 items/hour=481,174
05:03:58 INF Processing BHL items items=180,000 items/hour=501,707
05:04:26 INF Processing BHL items items=190,000 items/hour=518,120
05:05:07 INF Processing BHL items items=200,000 items/hour=529,207
05:06:01 INF Processing BHL items items=210,000 items/hour=534,408
05:07:01 INF Processing BHL items items=220,000 items/hour=537,047
05:08:18 INF Processing BHL items items=230,000 items/hour=533,448
```
