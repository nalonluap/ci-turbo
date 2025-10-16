CREATE TABLE IF NOT EXISTS default.metrics (
                                               ProjectID    String,
                                               PipelineID   String,
                                               JobID        String,
                                               JobName      String,
                                               RunnerID     String,
                                               StartTime    DateTime,
                                               DurationMs   Int64,
                                               ExitCode     Int32,
                                               CacheHit     UInt8,
                                               QueueTimeMs  Int64
) ENGINE = MergeTree()
    PARTITION BY toYYYYMM(StartTime)
    ORDER BY (ProjectID, StartTime);