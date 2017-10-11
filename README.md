# gobench
 Sysbench for  mysql create bench html
 
    test.conf 
    [DEFAULT]
    time_step =  2
    thread_list = 6, 10,15,20,30,40,50,70,100,120,180,200,300,400,500,800,1000
    db-driver=mysql 
    time=30



    [server1]
    mysql-host = 192.168.x.x
    mysql-port = 3306
    mysql-user = xxxxx
    mysql-password=xxxxx
    mysql-db=sbtest 
    cmd1 = "sysbench   --mysql-ignore-errors=all   --skip_trx=on   --db-ps-mode=disable  /usr/share/sysbench/oltp_read_only.lua  run"


 ./gobench  -f   test.conf 
 
 
