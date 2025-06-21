#!/bin/bash


log() {
	echo "$(date '+%Y-%m-%d %H:%M:%S') $*" | tee -a "$LOGFILE"
}

ENVFILE=approval.env
if [ ! -s "$ENVFILE" ];then
	log "envファイルが存在しません"	
	exit 1
fi
source "$ENVFILE"

log "ログイン開始...."
curl -c "$COOKEITEMPFILE" -X POST "${URL}${LOGINURL}"  -d "loginId=$LOGINID" -d "password=$PASSWORD" -d 'roll=staff' | tee -a "$LOGFILE"
if [ ! -s "$COOKEITEMPFILE" ]; then
	log "Cookie保存失敗..."
	exit 1
fi

log "経費承認開始...."
today=$(date +%Y%m%d)
result=$(curl -b "$COOKEITEMPFILE" -X POST "${URL}${EXPENSEURL}" | tee -a "$LOGFILE")
if echo "$result" | grep -q "$ERROR"; then
	log "経費承認失敗..."
	rm -f "$COOKEITEMPFILE"
	exit 1
fi
	

log "勤務表承認開始...."
result=$(curl -b "$COOKEITEMPFILE" -X POST "${URL}${ROSTERURL}" | tee -a "$LOGFILE")
if echo "$result" | grep -q "$ERROR"; then
	log "勤務表承認失敗..."
	rm -f "$COOKEITEMPFILE"
	exit 1
fi

rm -f "$COOKEITEMPFILE"
log "完了..."




