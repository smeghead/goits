VERSION = 0.0.01
COPYRIGHT = Copyright\ smeghead\ 2013\ -\ 2013

SRCS = **/*.go
PROGRAM = multivac

.PHONY: $(PROGRAM) clean run mergeresource compileresource

$(PROGRAM): compileresource
	go build multivac.go

clean:
	rm -f $(PROGRAM)
	rm -rf ./locale/ja
	rm -rf ./locale/en
	rm -rf ./locale/zh
	rm -rf ./locale/zh_TW

run: $(PROGRAM)
	./$(PROGRAM)

mergeresource:
	xgettext --from-code=utf-8 -k_ --msgid-bugs-address=smeghead@users.sourceforge.jp -L C -p locale  *.c js/*.js
	msgmerge -U locale/ja.po locale/messages.po
	msgmerge -U locale/en.po locale/messages.po
	msgmerge -U locale/zh.po locale/messages.po
	msgmerge -U locale/zh_TW.po locale/messages.po

compileresource: locale/ja/LC_MESSAGES/$(PROGRAM).mo locale/en/LC_MESSAGES/$(PROGRAM).mo locale/zh/LC_MESSAGES/$(PROGRAM).mo locale/zh_TW/LC_MESSAGES/$(PROGRAM).mo 

locale/ja/LC_MESSAGES/$(PROGRAM).mo: locale/ja.po 
	mkdir -p locale/ja/LC_MESSAGES
	msgfmt -o locale/ja/LC_MESSAGES/$(PROGRAM).mo locale/ja.po 

locale/en/LC_MESSAGES/$(PROGRAM).mo: locale/en.po 
	mkdir -p locale/en/LC_MESSAGES
	msgfmt -o locale/en/LC_MESSAGES/$(PROGRAM).mo locale/en.po 

locale/zh/LC_MESSAGES/$(PROGRAM).mo: locale/zh.po 
	mkdir -p locale/zh/LC_MESSAGES
	msgfmt -o locale/zh/LC_MESSAGES/$(PROGRAM).mo locale/zh.po 

locale/zh_TW/LC_MESSAGES/$(PROGRAM).mo: locale/zh_TW.po 
	mkdir -p locale/zh_TW/LC_MESSAGES
	msgfmt -o locale/zh_TW/LC_MESSAGES/$(PROGRAM).mo locale/zh_TW.po 
