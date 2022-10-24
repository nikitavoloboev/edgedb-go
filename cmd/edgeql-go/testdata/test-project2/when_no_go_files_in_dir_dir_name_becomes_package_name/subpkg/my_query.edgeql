create scalar type MyScalar extending int64;
create scalar type MyEnum extending enum<This, That>;

select {
	a := <uuid>$a,
	b := <optional uuid>$b,
	c := <str>$c,
	d := <optional str>$d,
	e := <bytes>$e,
	f := <optional bytes>$f,
	g := <int16>$g,
	h := <optional int16>$h,
	i := <int32>$i,
	j := <optional int32>$j,
	k := <int64>$k,
	l := <optional int64>$l,
	m := <float32>$m,
	n := <optional float32>$n,
	o := <float64>$o,
	p := <optional float64>$p,
	q := <bool>$q,
	r := <optional bool>$r,
	s := <datetime>$s,
	t := <optional datetime>$t,
	u := <cal::local_datetime>$u,
	v := <optional cal::local_datetime>$v,
	w := <cal::local_date>$w,
	x := <optional cal::local_date>$x,
	y := <cal::local_time>$y,
	z := <optional cal::local_time>$z,
	aa := <duration>$aa,
	ab := <optional duration>$ab,
	ac := <bigint>$ac,
	ad := <optional bigint>$ad,
	ae := <cal::relative_duration>$ae,
	af := <optional cal::relative_duration>$af,
	ag := <cal::date_duration>$ag,
	ah := <optional cal::date_duration>$ah,
	ai := <cfg::memory>$ai,
	aj := <optional cfg::memory>$aj,
	ak := <range<int32>>$ak,
	al := <optional range<int32>>$al,
	am := <range<int64>>$am,
	an := <optional range<int64>>$an,
	ao := <range<float32>>$ao,
	ap := <optional range<float32>>$ap,
	aq := <range<float64>>$aq,
	ar := <optional range<float64>>$ar,
	as := <range<datetime>>$as,
	at := <optional range<datetime>>$at,
	au := <range<cal::local_datetime>>$au,
	av := <optional range<cal::local_datetime>>$av,
	aw := <range<cal::local_date>>$aw,
	ax := <optional range<cal::local_date>>$ax,
	ay := <MyScalar>1,
	az := <optional MyScalar>{},
	ba := MyEnum.This,
	bb := <optional MyEnum>{},
}
