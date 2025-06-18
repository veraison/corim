SRCS := $(wildcard *.diag)
CBOR := $(SRCS:.diag=.cbor)
B64U := $(SRCS:.diag=.b64u)
HEX := $(SRCS:.diag=.hex)

%.cbor: %.diag ; diag2cbor.rb < $< > $@

CLEANFILES += $(CBOR)

%.b64u: %.cbor ; basenc --base64url --wrap 0 $< | tr -d '=' > $@

CLEANFILES += $(B64U)

%.hex: %.cbor ; xxd -p -c0 $< > $@

CLEANFILES += $(HEX)

all: $(CBOR) $(B64U) $(HEX)

clean: ; $(RM) -f $(CLEANFILES)
