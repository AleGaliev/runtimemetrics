```pprof
File: main
Type: cpu
Time: 2025-11-25 12:08:03 MSK
Duration: 240.04s, Total samples = 580ms ( 0.24%)
Showing nodes accounting for -70ms, 12.07% of 580ms total
Dropped 1 node (cum <= 2.90ms)
flat  flat%   sum%        cum   cum%
-40ms  6.90%  6.90%      -40ms  6.90%  runtime.pthread_cond_signal
30ms  5.17%  1.72%       30ms  5.17%  runtime.kevent
-30ms  5.17%  6.90%      -30ms  5.17%  runtime.typePointers.next
20ms  3.45%  3.45%       20ms  3.45%  runtime.pthread_cond_wait
-20ms  3.45%  6.90%      -20ms  3.45%  runtime.pthread_kill
-20ms  3.45% 10.34%      -20ms  3.45%  syscall.syscall
-10ms  1.72% 12.07%      -10ms  1.72%  runtime.(*mLockProfile).recordUnlock
-10ms  1.72% 13.79%      -10ms  1.72%  runtime.(*sweepLocked).sweep
-10ms  1.72% 15.52%      -10ms  1.72%  runtime.findObject
10ms  1.72% 13.79%       10ms  1.72%  runtime.madvise
10ms  1.72% 12.07%       10ms  1.72%  runtime.usleep
0     0% 12.07%       10ms  1.72%  bufio.(*Reader).ReadLine
0     0% 12.07%       10ms  1.72%  bufio.(*Reader).ReadSlice
0     0% 12.07%       10ms  1.72%  bufio.(*Reader).fill
0     0% 12.07%      -10ms  1.72%  bufio.(*Writer).Flush
0     0% 12.07%      -10ms  1.72%  database/sql.(*DB).Begin (inline)
0     0% 12.07%      -10ms  1.72%  database/sql.(*DB).BeginTx
0     0% 12.07%      -10ms  1.72%  database/sql.(*DB).BeginTx.func1
0     0% 12.07%       10ms  1.72%  database/sql.(*DB).QueryContext
0     0% 12.07%       10ms  1.72%  database/sql.(*DB).QueryContext.func1
0     0% 12.07%       10ms  1.72%  database/sql.(*DB).QueryRowContext (inline)
0     0% 12.07%      -10ms  1.72%  database/sql.(*DB).begin
0     0% 12.07%      -10ms  1.72%  database/sql.(*DB).beginDC
0     0% 12.07%      -10ms  1.72%  database/sql.(*DB).beginDC.func1
0     0% 12.07%      -10ms  1.72%  database/sql.(*DB).prepareDC
0     0% 12.07%      -10ms  1.72%  database/sql.(*DB).prepareDC.func2
0     0% 12.07%       10ms  1.72%  database/sql.(*DB).query
0     0% 12.07%       10ms  1.72%  database/sql.(*DB).queryDC
0     0% 12.07%       10ms  1.72%  database/sql.(*DB).queryDC.func1
0     0% 12.07%      -50ms  8.62%  database/sql.(*DB).retry
0     0% 12.07%      -50ms  8.62%  database/sql.(*Stmt).ExecContext
0     0% 12.07%      -50ms  8.62%  database/sql.(*Stmt).ExecContext.func1
0     0% 12.07%      -10ms  1.72%  database/sql.(*Tx).Prepare (inline)
0     0% 12.07%      -10ms  1.72%  database/sql.(*Tx).PrepareContext
0     0% 12.07%      -10ms  1.72%  database/sql.(*driverConn).prepareLocked
0     0% 12.07%      -10ms  1.72%  database/sql.ctxDriverBegin
0     0% 12.07%      -10ms  1.72%  database/sql.ctxDriverPrepare
0     0% 12.07%       10ms  1.72%  database/sql.ctxDriverQuery
0     0% 12.07%      -50ms  8.62%  database/sql.ctxDriverStmtExec
0     0% 12.07%      -50ms  8.62%  database/sql.resultFromStatement
0     0% 12.07%      -10ms  1.72%  database/sql.withLock
0     0% 12.07%      -10ms  1.72%  encoding/json.(*Decoder).Decode
0     0% 12.07%      -10ms  1.72%  encoding/json.(*Decoder).readValue
0     0% 12.07%      -10ms  1.72%  encoding/json.stateInString
0     0% 12.07%      -10ms  1.72%  gcWriteBarrier
0     0% 12.07%      -40ms  6.90%  github.com/AleGaliev/runtimemetrics/internal/handler.CreateMyHandler.GzipMiddlewareHandler.func6.1
0     0% 12.07%      -40ms  6.90%  github.com/AleGaliev/runtimemetrics/internal/handler.CreateMyHandler.MetricValidateMiddleware.func7.1
0     0% 12.07%      -40ms  6.90%  github.com/AleGaliev/runtimemetrics/internal/handler.CreateMyHandler.MiddlewareHandlerLogger.func5.1
0     0% 12.07%      -30ms  5.17%  github.com/AleGaliev/runtimemetrics/internal/handler.CreateMyHandler.func4.AuditMiddleware.1.1
0     0% 12.07%      -60ms 10.34%  github.com/AleGaliev/runtimemetrics/internal/handler.MyHandler.ServeHTTPBatchUpdate
0     0% 12.07%       30ms  5.17%  github.com/AleGaliev/runtimemetrics/internal/observer.(*Event).Notify
0     0% 12.07%       30ms  5.17%  github.com/AleGaliev/runtimemetrics/internal/repository.(*AuditSaveFile).SendAudit
0     0% 12.07%      -60ms 10.34%  github.com/AleGaliev/runtimemetrics/internal/storage.(*PostgresDBStorage).BatchUpdateMetrics
0     0% 12.07%       10ms  1.72%  github.com/AleGaliev/runtimemetrics/internal/storage.(*PostgresDBStorage).counterManipulation
0     0% 12.07%      -40ms  6.90%  github.com/go-chi/chi/v5.(*ChainHandler).ServeHTTP
0     0% 12.07%      -30ms  5.17%  github.com/go-chi/chi/v5.(*Mux).Mount.func1
0     0% 12.07%      -40ms  6.90%  github.com/go-chi/chi/v5.(*Mux).ServeHTTP
0     0% 12.07%      -40ms  6.90%  github.com/go-chi/chi/v5.(*Mux).routeHTTP
0     0% 12.07%      -10ms  1.72%  github.com/jackc/pgx/v5.(*Conn).BeginTx
0     0% 12.07%      -60ms 10.34%  github.com/jackc/pgx/v5.(*Conn).Exec
0     0% 12.07%      -10ms  1.72%  github.com/jackc/pgx/v5.(*Conn).Prepare
0     0% 12.07%       10ms  1.72%  github.com/jackc/pgx/v5.(*Conn).Query
0     0% 12.07%      -60ms 10.34%  github.com/jackc/pgx/v5.(*Conn).exec
0     0% 12.07%      -50ms  8.62%  github.com/jackc/pgx/v5.(*Conn).execPrepared
0     0% 12.07%      -10ms  1.72%  github.com/jackc/pgx/v5.(*Conn).execSimpleProtocol
0     0% 12.07%      -10ms  1.72%  github.com/jackc/pgx/v5/pgconn.(*MultiResultReader).NextResult
0     0% 12.07%      -10ms  1.72%  github.com/jackc/pgx/v5/pgconn.(*MultiResultReader).receiveMessage
0     0% 12.07%      -40ms  6.90%  github.com/jackc/pgx/v5/pgconn.(*PgConn).ExecPrepared
0     0% 12.07%      -10ms  1.72%  github.com/jackc/pgx/v5/pgconn.(*PgConn).Prepare
0     0% 12.07%      -40ms  6.90%  github.com/jackc/pgx/v5/pgconn.(*PgConn).execExtendedSuffix
0     0% 12.07%      -30ms  5.17%  github.com/jackc/pgx/v5/pgconn.(*PgConn).flushWithPotentialWriteReadDeadlock
0     0% 12.07%      -30ms  5.17%  github.com/jackc/pgx/v5/pgconn.(*PgConn).peekMessage
0     0% 12.07%      -20ms  3.45%  github.com/jackc/pgx/v5/pgconn.(*PgConn).receiveMessage
0     0% 12.07%      -10ms  1.72%  github.com/jackc/pgx/v5/pgconn.(*ResultReader).readUntilRowDescription
0     0% 12.07%      -30ms  5.17%  github.com/jackc/pgx/v5/pgconn/internal/bgreader.(*BGReader).Read
0     0% 12.07%      -30ms  5.17%  github.com/jackc/pgx/v5/pgproto3.(*Frontend).Flush
0     0% 12.07%      -30ms  5.17%  github.com/jackc/pgx/v5/pgproto3.(*Frontend).Receive
0     0% 12.07%      -30ms  5.17%  github.com/jackc/pgx/v5/pgproto3.(*chunkReader).Next
0     0% 12.07%      -10ms  1.72%  github.com/jackc/pgx/v5/stdlib.(*Conn).BeginTx
0     0% 12.07%      -50ms  8.62%  github.com/jackc/pgx/v5/stdlib.(*Conn).ExecContext
0     0% 12.07%      -10ms  1.72%  github.com/jackc/pgx/v5/stdlib.(*Conn).PrepareContext
0     0% 12.07%       10ms  1.72%  github.com/jackc/pgx/v5/stdlib.(*Conn).QueryContext
0     0% 12.07%      -50ms  8.62%  github.com/jackc/pgx/v5/stdlib.(*Stmt).ExecContext
0     0% 12.07%       10ms  1.72%  internal/poll.(*FD).Accept
0     0% 12.07%      -20ms  3.45%  internal/poll.(*FD).Read
0     0% 12.07%      -40ms  6.90%  internal/poll.(*FD).Write
0     0% 12.07%       10ms  1.72%  internal/poll.accept
0     0% 12.07%      -60ms 10.34%  internal/poll.ignoringEINTRIO (inline)
0     0% 12.07%      -30ms  5.17%  io.ReadAtLeast
0     0% 12.07%       10ms  1.72%  main.main
0     0% 12.07%       10ms  1.72%  net.(*TCPListener).Accept
0     0% 12.07%       10ms  1.72%  net.(*TCPListener).accept
0     0% 12.07%      -20ms  3.45%  net.(*conn).Read
0     0% 12.07%      -40ms  6.90%  net.(*conn).Write
0     0% 12.07%      -20ms  3.45%  net.(*netFD).Read
0     0% 12.07%      -40ms  6.90%  net.(*netFD).Write
0     0% 12.07%       10ms  1.72%  net.(*netFD).accept
0     0% 12.07%       10ms  1.72%  net/http.(*Server).ListenAndServe
0     0% 12.07%       10ms  1.72%  net/http.(*Server).Serve
0     0% 12.07%       10ms  1.72%  net/http.(*conn).readRequest
0     0% 12.07%      -40ms  6.90%  net/http.(*conn).serve
0     0% 12.07%       10ms  1.72%  net/http.(*connReader).Read
0     0% 12.07%      -10ms  1.72%  net/http.(*response).finishRequest
0     0% 12.07%      -40ms  6.90%  net/http.HandlerFunc.ServeHTTP
0     0% 12.07%       10ms  1.72%  net/http.ListenAndServe (inline)
0     0% 12.07%      -10ms  1.72%  net/http.checkConnErrorWriter.Write
0     0% 12.07%       10ms  1.72%  net/http.readRequest
0     0% 12.07%      -40ms  6.90%  net/http.serverHandler.ServeHTTP
0     0% 12.07%       10ms  1.72%  net/textproto.(*Reader).ReadLine (inline)
0     0% 12.07%       10ms  1.72%  net/textproto.(*Reader).readLineSlice
0     0% 12.07%       30ms  5.17%  os.OpenFile
0     0% 12.07%       30ms  5.17%  os.ignoringEINTR (inline)
0     0% 12.07%       10ms  1.72%  os.newFile
0     0% 12.07%       10ms  1.72%  os.newFile.func1 (inline)
0     0% 12.07%       20ms  3.45%  os.open (inline)
0     0% 12.07%       30ms  5.17%  os.openFileNolog
0     0% 12.07%       20ms  3.45%  os.openFileNolog.func1 (inline)
0     0% 12.07%      -20ms  3.45%  runtime.(*gcControllerState).enlistWorker
0     0% 12.07%      -20ms  3.45%  runtime.(*gcWork).balance
0     0% 12.07%      -10ms  1.72%  runtime.(*mcache).prepareForSweep
0     0% 12.07%      -10ms  1.72%  runtime.(*mcache).releaseAll
0     0% 12.07%      -10ms  1.72%  runtime.(*mcentral).uncacheSpan
0     0% 12.07%       10ms  1.72%  runtime.(*mheap).alloc.func1
0     0% 12.07%       10ms  1.72%  runtime.(*mheap).allocSpan
0     0% 12.07%      -40ms  6.90%  runtime.entersyscall_sysmon
0     0% 12.07%       30ms  5.17%  runtime.findRunnable
0     0% 12.07%      -10ms  1.72%  runtime.forEachP (inline)
0     0% 12.07%      -10ms  1.72%  runtime.gcBgMarkWorker
0     0% 12.07%      -20ms  3.45%  runtime.gcBgMarkWorker.func2
0     0% 12.07%      -20ms  3.45%  runtime.gcDrain
0     0% 12.07%      -20ms  3.45%  runtime.gcDrainMarkWorkerDedicated (inline)
0     0% 12.07%      -10ms  1.72%  runtime.gcMarkDone
0     0% 12.07%      -10ms  1.72%  runtime.gcMarkTermination
0     0% 12.07%       10ms  1.72%  runtime.gcMarkTermination.func3
0     0% 12.07%       10ms  1.72%  runtime.gcMarkTermination.func4.1
0     0% 12.07%       10ms  1.72%  runtime.goexit0
0     0% 12.07%       10ms  1.72%  runtime.gopreempt_m (inline)
0     0% 12.07%       10ms  1.72%  runtime.goschedImpl
0     0% 12.07%       20ms  3.45%  runtime.mPark (inline)
0     0% 12.07%       10ms  1.72%  runtime.main
0     0% 12.07%       20ms  3.45%  runtime.mcall
0     0% 12.07%       10ms  1.72%  runtime.morestack
0     0% 12.07%       40ms  6.90%  runtime.netpoll
0     0% 12.07%      -10ms  1.72%  runtime.netpollBreak (inline)
0     0% 12.07%       20ms  3.45%  runtime.netpollgoready.goready.func1
0     0% 12.07%      -10ms  1.72%  runtime.newproc.func1
0     0% 12.07%       10ms  1.72%  runtime.newstack
0     0% 12.07%       20ms  3.45%  runtime.notesleep
0     0% 12.07%      -40ms  6.90%  runtime.notewakeup
0     0% 12.07%       10ms  1.72%  runtime.park_m
0     0% 12.07%      -20ms  3.45%  runtime.preemptM
0     0% 12.07%      -20ms  3.45%  runtime.preemptone
0     0% 12.07%       10ms  1.72%  runtime.ready
0     0% 12.07%       10ms  1.72%  runtime.resetspinning
0     0% 12.07%       10ms  1.72%  runtime.runqgrab
0     0% 12.07%       10ms  1.72%  runtime.runqsteal
0     0% 12.07%       30ms  5.17%  runtime.schedule
0     0% 12.07%       20ms  3.45%  runtime.semasleep
0     0% 12.07%      -40ms  6.90%  runtime.semawakeup
0     0% 12.07%      -10ms  1.72%  runtime.send.goready.func1
0     0% 12.07%      -20ms  3.45%  runtime.signalM (inline)
0     0% 12.07%       10ms  1.72%  runtime.startTheWorldWithSema
0     0% 12.07%       10ms  1.72%  runtime.stealWork
0     0% 12.07%       10ms  1.72%  runtime.stopm
0     0% 12.07%       10ms  1.72%  runtime.sysUsed (inline)
0     0% 12.07%       10ms  1.72%  runtime.sysUsedOS (inline)
0     0% 12.07%      -50ms  8.62%  runtime.systemstack
0     0% 12.07%      -10ms  1.72%  runtime.unlock (inline)
0     0% 12.07%      -10ms  1.72%  runtime.unlock2
0     0% 12.07%      -10ms  1.72%  runtime.unlockWithRank (inline)
0     0% 12.07%      -10ms  1.72%  runtime.wakeNetpoll
0     0% 12.07%      -10ms  1.72%  runtime.wbBufFlush
0     0% 12.07%      -10ms  1.72%  runtime.wbBufFlush.func1
0     0% 12.07%      -10ms  1.72%  runtime.wbBufFlush1
0     0% 12.07%       10ms  1.72%  syscall.CloseOnExec (inline)
0     0% 12.07%       10ms  1.72%  syscall.Fstat
0     0% 12.07%       20ms  3.45%  syscall.Open
0     0% 12.07%      -20ms  3.45%  syscall.Read (inline)
0     0% 12.07%      -40ms  6.90%  syscall.Write (inline)
0     0% 12.07%       10ms  1.72%  syscall.fcntl
0     0% 12.07%      -20ms  3.45%  syscall.read
0     0% 12.07%      -40ms  6.90%  syscall.write
```