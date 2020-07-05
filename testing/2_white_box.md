1. the original function to be tested is in storage1. in storage2 we rewrite that function by extracting its email stuff into another function. so storage1 and storage2 do the same thing.

2. then we change that email function (like stub in sinon) with another function in the test. and then we restore it back.

3. seems that function in storage2 will be used in production instead. so in whitebox testing, we also change the functions inside production code.