#!/usr/bin/env perl
use strict;
use warnings;
use Getopt::Long;

my $usage = q{
#> Usage: fork [-d] -p=PIDS_FILE -n=N [-k] [-ka] [-r] [-dela=D] "COMMAND"
#>
#> Запускает N процессов COMMAND, записывая номера процессов в файл PIDS_FILE.
#> -d|deb|debug - выводит отладочные сообщения в STDOUT
#> -p|pf|pidfile|pidsfile=PIDS_FILE - путь и имя файла. Если отсутствует - будет создан.
#> -n|number=N - требуемое количество процессов. 
#> 	При старте будет считан PIDS_FILE с номерами процессов. 
#> 	Они будут проверены на существование и недостающее количество процессов (или нисколько) будет запущено.
#>	Если количество фактически работающих процессов больше требуемого - то,
#>		если указана опция -k|kill, то лишние - будут прерваны отправкой kill.
#>
#> -single - тоже что и "-n=1", используется для заранее однопотоковых скриптов, в целях предотвратить случайное изменение парамера -n на большее число.
#> -ka|killall|kilall|kilal|kila - ВСЕМ процессам из PIDS_FILE посылается kill.	Кроме них - еще и всем их дочерним процессам.
#> -wait - ждать завершения запущенных в этот раз дочерних процессов.
#> -status - возвращать статус возврата от вызова wait. Включает --wait. 
#> -r|rep|repo|report - только вывод отчета, о текущем состоянии. 
#> -lsof - выводит открытые файлы (на момент запуска).
#> -dela|delay=D, где D - количество минут. 
#> 	Второй и каждый последующий процесс будет предваряться командой sleep с увеличивающейся задержкой в минутах.
#> --comment="comment text" [ -comment="comment text2" ] - комментарий. Никак не используется. Применяется лишь для комментирования кода.
#> --decrease - уменьшить на один кол-во работающих процессов
#>
#> Samples:
#> fork -pf=ECHO.PIDS -n=4 -dela=1 "echo AAA" # - запустить 4 процесса, печатающих AAA, первый - сейчас, 2-й - через 1 минуту,..,4-й - через 3 мин.
#>  если повторить вышеприведенную команду менее чем через минуту (т.е. когда останутся 3 процесса) - запустится 1 (недостающий до 4-х) процесс.
#> fork -pf=ECHO.PIDS -repo  # какие процессы из записанных в ECHO.PIDS реально работают. Обновляет записи в ECHO.PIDS после проверки.
#> fork -pf=ECHO.PIDS -n=2 -dela=1 -k "echo AAA" # - если работает больше 2 процессов, то остальным послать kill. Если меньше - запустить недостающие.
#> Убивание родительского процесса не всегда убивает дочерние, поэтому убиваются еще и все дочерние.
#> fork -pf=ECHO.PIDS -kila # послать всем - kill.

};

$|=1;
my $pf;# = "TEST.PIDS";
my $deb;
my $lsof;
my $kill;
my $killall;
my $need_cnt;
my $single;
my $report_only;
my $delay;
my $help;
my $wait;
my $getstatus;
my $ed;
my @comments;
my $decrease;
GetOptions(
    "p|pf|pidfile|pidsfile=s" => \$pf,
    "d|deb|debug" => \$deb,
    "k|kill" => \$kill,
    "n|number=i" => \$need_cnt,
    "single" => \$single,
    "ka|killall|kilall|kilal|kila" => \$killall,
    "r|rep|repo|report" => \$report_only,
    "wait" => \$wait,
    "status" => \$getstatus,
    "lsof" => \$lsof,
    "dela|delay=i" => \$delay,
    "h|help" => \$help,
    "comment=s" => \@comments,
    "ed=s" => \$ed,
    "decrease" => \$decrease, 
) or die "Bad options!";
die $usage if $help;
$pf or die "pidsfile!";

my $cmd = shift;
if ($decrease){
 $cmd ||= "true";
 $kill = 1;
}

if ( $killall ){
    ($need_cnt, $single, $kill) = (0,0,1);
    $cmd ||= "true";
}

if ( $report_only ){
    ( $deb, $kill, $killall, $need_cnt, $single ) = (1, 0, 0, 0, 0);
    $cmd ||= "true";
}

$need_cnt||=1 if $single;
defined $need_cnt or die "need_cnt!";

if ( -e $pf ){
    die "'$pf' - is executable!" if -x $pf;
    die "'$pf' - have not write perm!" if !-w $pf;
    die "'$pf' - have not read perm!" if !-r $pf;
}else{
    my $pfdir = `dirname $pf`; chomp $pfdir;
    die "'$pfdir' - have not write perm!" if !-w $pfdir;
}

$cmd or die "cmd!";

my %eds = (s=>"sec", m=>"min", h=>"hour");
$ed||="m";
$eds{ $ed } or die "Bad ed! m.b.".join(", ",keys %eds);

my ($writed_pids, $runned_pids) = pids( $pf );

$wait=1 if $getstatus;
$SIG{CHLD}="IGNORE" if not $wait;

if ($decrease){
 # если попросили уменьшить, то просим на один работающий просесс меньше, чем есть 
 $need_cnt = @$runned_pids - 1; 
 $need_cnt = 0 if $need_cnt < 0;
}

if ( @$runned_pids < $need_cnt ){
    # если запущено меньше, чем нужно:
    my $n;
    for my $i (1..$need_cnt-@$runned_pids){
	  $n = ($i-1)*$delay if $delay;
	  if (my $pid=fork){
	    print "forked $pid\n" if $deb;
	    `echo $pid >> $pf`;
	  }else{
	    my $sleep_cmd='';
	    if ( $delay and $n ){
		$sleep_cmd = "sleep $n"."$ed;";
	    }
	    exec "$sleep_cmd $cmd";
	  }
    }
    if ($wait){
        wait;
        if ( $getstatus ){
          #warn "?: $?";
          if ($? == -1) {
            warn "I havent childs" if $deb;
            exit $?;
          }
          elsif ( $? & 127 ) {
            warn "child died" if $deb;
            exit ($? & 127);
          }
          else {
            warn "child exited" if $deb;
            exit ($? >> 8);
          }
        }
    }    
}
elsif( @$runned_pids > $need_cnt and $kill ){
    for ( 1..@$runned_pids-$need_cnt ){
	my $pid = @$runned_pids[$_-1];
	my @childs = childs( $pid );
	print "kill $pid\n" if $deb;
	kill("TERM", $pid);
	for my $ch ( @childs ){
	    print "\tkill child process $ch\n" if $deb;
	    kill "TERM", $ch;
	}
    }
    pids( $pf );
}

#--------------------------------subs-------------------
sub pids{
 my $pf = shift or die "pidfile!";
 my @writed_pids = split(/\s+/, `cat $pf`);
 print "writed pids: @writed_pids\n" if $deb;
 my @runned_pids = grep { $_ && kill(0, $_) } @writed_pids;
 print "runned pids: @runned_pids\n" if $deb;
 my @childs = childs( @runned_pids );
 print "child process: @childs\n" if @childs and $deb;
 if ($lsof){
    print "Opened files: ". lsof( @runned_pids, @childs )."\n";
 } 
 `echo -n '' > $pf`;
 `echo $_ >> $pf` for @runned_pids;
 return ( \@writed_pids, \@runned_pids );
}

sub childs{
 my @pids = @_;
 my @rv;
 return if !@pids;
 for my $pid (@pids){
    my $childs = `pgrep -P $pid 2>/dev/null`;
    my @childs = split /\n/, $childs;
    #warn "Childs for $pid: @childs" if $deb;
    push @rv, @childs;
 }
 return if !@rv;
 return ( @rv, childs(@rv) );
}

use Data::Dumper;

sub lsof{
 my @text = split /\n/, `lsof -Fpnat`;
 my (%h, $p,$current_p );
 for ( @text ){
    if ( ($p) = /^p(\d+)/ ){
	$current_p = $p;
    }
    if ( my ($n) = /^n(.+)/ ){
	$h{$current_p}{n} = $n;
    }
    if ( my ($ac) = /^a(.+)/ ){
	$h{$current_p}{a} = $ac;
    }
    if ( my ($t) = /^t(\w+)/ ){
	$h{$current_p}{t} = $t;
    }
 }
 my %rv;
 for (@_){
    my $el = $h{ $_ }||{};
    $el->{t}||"" eq "REG" or next;
    #$el->{n}||"" eq "pipe" and next;
    $rv{ $_ } = $el;
 }
 #die Dumper \%rv;
 return wantarray ? %rv : join ', ', map { "$_:". join(' ', grep {$_} @{$rv{$_}}{"n","a"}) } keys %rv;    
}

