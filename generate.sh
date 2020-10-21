#!/bin/sh

PID=/tmp/app.pid
LOG=/tmp/app.log

start() {
  if [ -f $PID ]
  then
    echo
    echo "Already started. PID [$( cat $PID )]"
  else
    touch $PID

    if [ ! -d testmedia ]
    then
      mkdir testmedia
    fi

    rm -rf testmedia/live
    mkdir testmedia/live
    mkdir testmedia/live/test
    mkdir testmedia/live/test/2M
    mkdir testmedia/live/test/4M
    mkdir testmedia/live/test/audio
    
    if nohup mediastreamsegmenter -w 202 -t 4 224.0.0.50:9120 -s 16 -D -T -i playlist.m3u8 --cmaf-fragmented --video-only -f testmedia/live/test/2M/ >>$LOG 2>&1 &
    then echo "START L"
    else echo "Error..."
      rm $PID
      return 1
    fi

    if nohup mediastreamsegmenter -w 202 -t 4 224.0.0.50:9121 -s 16 -D -T -i playlist.m3u8 --cmaf-fragmented --video-only -f testmedia/live/test/4M/ >>$LOG 2>&1 & 
    then echo "START O"
    else echo "Error..."
      rm $PID
      return 1
    fi

    if nohup mediastreamsegmenter -w 202 -t 4 224.0.0.50:9121 -s 16 -D -T -i playlist.m3u8 --cmaf-fragmented --audio-only -f testmedia/live/test/audio/ >>$LOG 2>&1 & 
    then echo "START A"
    else echo "Error..."
      rm $PID
      return 1
    fi

    if nohup tsrecompressor -L 224.0.0.50:9120 -O 224.0.0.50:9121 -h -g -x -a >>$LOG 2>&1 &
    then echo $! >$PID
      echo "Done."
    else echo "Error..."
      rm $PID
    fi
  fi
}

kill_cmd() {
  SIGNAL="";
  MSG="Killing "
  while true
  do
    LIST=`ps -ef | grep -v grep | grep -e mediastreamsegmenter -e tsrecompressor | awk '{print $2}'`
    if [ "$LIST" ]
    then
      echo; echo "$MSG $LIST" ; echo
      echo $LIST | xargs kill $SIGNAL
      sleep 2
      SIGNAL="-9" ; MSG="Killing $SIGNAL"
      if [ -f $PID ]
      then
        rm $PID
      fi
    else
      echo; echo "All killed..." ; echo
      break
    fi
  done
}

stop() {
  if [ -f $PID ]
  then
    if kill $( cat $PID )
    then echo "Done."
    fi
    rm $PID
    kill_cmd
  else
    echo "No pid file. Already stopped"
  fi
}

case "$1" in
  'start')
    start
    ;;
  'stop')
    stop
    ;;
esac

exit 0
