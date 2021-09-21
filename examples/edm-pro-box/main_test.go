package main

import (
	"strings"
	"testing"
)

func TestGCMC(t *testing.T) {
	g := gcmc()
	got := strings.Split(g.String(), "\n")
	want := strings.Split(gcmcOut, "\n")

	if len(got) != len(want) {
		t.Errorf("gcmc = %v lines, want %v", len(got), len(want))
	}

	for i, line := range got {
		if i == 1 {
			continue
		}
		if i < len(want) && line != want[i] {
			t.Fatalf("gcmc line #%v:\n%v\nwant:\n%v", i+1, line, want[i])
		}
	}
}

var gcmcOut = `(go-gcode compiled code, do not change)
(2021-09-11 08:33:05)
(-- prologue begin --)
G17 ( Use XY plane )
G21 ( Use mm )
G40 ( Cancel cutter radius compensation )
G49 ( Cancel tool length compensation )
G54 ( Default coordinate system )
G80 ( Cancel canned cycle )
G90 ( Use absolute distance mode )
G94 ( Units Per Minute feed rate mode )
G64 ( Enable path blending for best speed )
(-- prologue end --)
F500.00000000
G0 X-0.05000000 Y78.00000000
G2 X-0.05000000 Y78.00000000 I1.55000000 J0.00000000
G0 X17.87500000 Y20.00000000
G2 X17.87500000 Y20.00000000 I2.12500000 J0.00000000
G0 X53.45000000 Y1.50000000
G2 X53.45000000 Y1.50000000 I1.55000000 J0.00000000
G0 X87.87500000 Y20.00000000
G2 X87.87500000 Y20.00000000 I2.12500000 J0.00000000
G0 X106.95000000 Y78.00000000
G2 X106.95000000 Y78.00000000 I1.55000000 J0.00000000
G0 X87.87500000 Y136.00000000
G2 X87.87500000 Y136.00000000 I2.12500000 J0.00000000
G0 X53.45000000 Y154.50000000
G2 X53.45000000 Y154.50000000 I1.55000000 J0.00000000
G0 X17.87500000 Y136.00000000
G2 X17.87500000 Y136.00000000 I2.12500000 J0.00000000
G0 X0.00000000 Y0.00000000
G1 X110.00000000
G1 Y156.00000000
G1 X0.00000000
G1 Y0.00000000
G0 X113.00000000 Y150.00000000
G1 Y76.50000000
G1 X117.00000000
G1 Y77.75000000
G1 X119.50000000
G1 Y76.50000000
G1 X121.50000000
G1 Y73.50000000
G1 X119.50000000
G1 Y72.25000000
G1 X117.00000000
G1 Y73.50000000
G1 X113.00000000
G1 Y0.00000000
G1 X134.72452355
G1 Y4.00000000
G1 X133.47452355
G1 Y6.50000000
G1 X134.72452355
G1 Y8.50000000
G1 X137.72452355
G1 Y6.50000000
G1 X138.97452355
G1 Y4.00000000
G1 X137.72452355
G1 Y0.00000000
G1 X159.44904711
G1 X153.09725760 Y36.02278837
G1 X149.15802658 Y35.32819566
G1 X149.37508681 Y34.09718597
G1 X146.91306742 Y33.66306552
G1 X146.69600720 Y34.89407522
G1 X144.72639170 Y34.54677886
G1 X144.20544716 Y37.50120212
G1 X146.17506267 Y37.84849847
G1 X145.95800245 Y39.07950817
G1 X148.42002183 Y39.51362861
G1 X148.63708205 Y38.28261892
G1 X152.57631306 Y38.97721163
G1 X139.87273404 Y111.02278837
G1 X135.93350303 Y110.32819566
G1 X136.15056325 Y109.09718597
G1 X133.68854387 Y108.66306552
G1 X133.47148365 Y109.89407522
G1 X131.50186814 Y109.54677886
G1 X130.98092361 Y112.50120212
G1 X132.95053912 Y112.84849847
G1 X132.73347889 Y114.07950817
G1 X135.19549828 Y114.51362861
G1 X135.41255850 Y113.28261892
G1 X139.35178951 Y113.97721163
G1 X133.00000000 Y150.00000000
G1 X124.50000000
G1 Y146.00000000
G1 X125.75000000
G1 Y143.50000000
G1 X124.50000000
G1 Y141.50000000
G1 X121.50000000
G1 Y143.50000000
G1 X120.25000000
G1 Y146.00000000
G1 X121.50000000
G1 Y150.00000000
G1 X113.00000000
G0 X164.00000000 Y0.00000000
G1 X151.23685894 Y72.38336985
G1 X155.17608995 Y73.07796256
G1 X155.39315018 Y71.84695287
G1 X157.85516956 Y72.28107331
G1 X157.63810934 Y73.51208300
G1 X159.60772484 Y73.85937936
G1 X159.08678031 Y76.81380262
G1 X157.11716480 Y76.46650626
G1 X156.90010458 Y77.69751595
G1 X154.43808520 Y77.26339551
G1 X154.65514542 Y76.03238582
G1 X150.71591441 Y75.33779311
G1 X137.95277335 Y147.72116295
G1 X159.34725258 Y151.49358688
G1 X160.04184529 Y147.55435587
G1 X158.81083559 Y147.33729564
G1 X159.24495604 Y144.87527626
G1 X160.47596573 Y145.09233648
G1 X160.82326209 Y143.12272098
G1 X163.77768534 Y143.64366551
G1 X163.43038899 Y145.61328102
G1 X164.66139868 Y145.83034124
G1 X164.22727824 Y148.29236062
G1 X162.99626855 Y148.07530040
G1 X162.30167583 Y152.01453141
G1 X183.69615506 Y155.78695534
G1 Y119.20845739
G1 X179.69615506
G1 Y120.45845739
G1 X177.19615506
G1 Y119.20845739
G1 X175.19615506
G1 Y116.20845739
G1 X177.19615506
G1 Y114.95845739
G1 X179.69615506
G1 Y116.20845739
G1 X183.69615506
G1 Y43.05146150
G1 X179.69615506
G1 Y44.30146150
G1 X177.19615506
G1 Y43.05146150
G1 X175.19615506
G1 Y40.05146150
G1 X177.19615506
G1 Y38.80146150
G1 X179.69615506
G1 Y40.05146150
G1 X183.69615506
G1 Y3.47296355
G1 X175.32528916 Y1.99695404
G1 X174.63069645 Y5.93618506
G1 X175.86170614 Y6.15324528
G1 X175.42758570 Y8.61526466
G1 X174.19657600 Y8.39820444
G1 X173.84927965 Y10.36781994
G1 X170.89485639 Y9.84687541
G1 X171.24215275 Y7.87725990
G1 X170.01114305 Y7.66019968
G1 X170.44526350 Y5.19818030
G1 X171.67627319 Y5.41524052
G1 X172.37086590 Y1.47600951
G1 X164.00000000 Y0.00000000
G0 X195.42850000 Y7.95950000
G2 X195.42850000 Y7.95950000 I1.55000000 J0.00000000
G0 X285.47150000
G2 X285.47150000 Y7.95950000 I1.55000000 J0.00000000
G0 X293.95000000 Y41.07849795
G2 X293.95000000 Y41.07849795 I1.55000000 J0.00000000
G0 Y117.23549384
G2 X293.95000000 Y117.23549384 I1.55000000 J0.00000000
G0 X285.47150000 Y148.04050000
G2 X285.47150000 Y148.04050000 I1.55000000 J0.00000000
G0 X195.42850000
G2 X195.42850000 Y148.04050000 I1.55000000 J0.00000000
G0 X186.95000000 Y117.23549384
G2 X186.95000000 Y117.23549384 I1.55000000 J0.00000000
G0 Y41.07849795
G2 X186.95000000 Y41.07849795 I1.55000000 J0.00000000
G0 X217.24750000 Y21.04050000
G2 X217.24750000 Y21.04050000 I4.75000000 J0.00000000
G0 X257.25250000
G2 X257.25250000 Y21.04050000 I4.75000000 J0.00000000
G0 X187.00000000 Y0.00000000
G1 X297.00000000
G1 Y158.40655145
G1 X187.00000000
G1 Y0.00000000
G0 X219.95000000 Y172.00000000
G2 X219.95000000 Y172.00000000 I1.55000000 J0.00000000
G0 X112.95000000
G2 X112.95000000 Y172.00000000 I1.55000000 J0.00000000
G0 X113.00000000 Y162.00000000
G1 X166.50000000
G1 Y166.00000000
G1 X165.25000000
G1 Y168.50000000
G1 X166.50000000
G1 Y170.50000000
G1 X169.50000000
G1 Y168.50000000
G1 X170.75000000
G1 Y166.00000000
G1 X169.50000000
G1 Y162.00000000
G1 X223.00000000
G1 Y181.47905547
G1 X113.00000000
G1 Y162.00000000
G0 X106.95000000 Y182.22452355
G2 X106.95000000 Y182.22452355 I1.55000000 J0.00000000
G0 X76.60000000 Y175.09961884
G1 Y180.05961884
G2 X88.40000000 Y180.05961884 I5.90000000 J-2.47991935
G1 Y175.09961884
G2 X76.60000000 Y175.09961884 I-5.90000000 J2.47991935
G0 X-0.05000000 Y182.22452355
G2 X-0.05000000 Y182.22452355 I1.55000000 J0.00000000
G0 X0.00000000 Y159.00000000
G1 X53.50000000
G1 Y163.00000000
G1 X52.25000000
G1 Y165.50000000
G1 X53.50000000
G1 Y167.50000000
G1 X56.50000000
G1 Y165.50000000
G1 X57.75000000
G1 Y163.00000000
G1 X56.50000000
G1 Y159.00000000
G1 X110.00000000
G1 Y205.44904711
G1 X0.00000000
G1 Y159.00000000
(-- epilogue begin --)
M30 (-- epilogue end --)
`