#!/bin/bash
printf "AVGPROC-%03d" $((RANDOM % 400))
cat <<EOF
% 0016
MDC READS-000001/SEC WRITES-000001/SEC HIT RATIO-100%
PAGING-1/SEC
Q0-00000 Q1-00000           Q2-00000 EXPAN-002 Q3-00033 EXPAN-002
EOF

