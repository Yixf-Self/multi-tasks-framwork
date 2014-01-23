#!/usr/bin/env perl
#
# This is a perl script illustrating how to use multi-threads and
# queue to do multi-tasks.
#
# It use the producer-consumer model.
#
#
# Author: Wei Shen <shenwei356#gmail.com>
# Site:   http://shenwei.me https://github.com/shenwei356
#

use strict;
use warnings;

use threads;
use Thread::Queue;

# number of wokers, the number of CPUs is recommended.
my $WORKER_NUM = 4;

# glob variable
my $queue;

MAIN: {
    die "Usage: $0 STRING [STRING ...]\n"
        unless @ARGV > 0;

    # create a queue
    $queue = Thread::Queue->new;

   # Put some thing into the queue
   # Situation
   # 1) If your to-do tasks could be defined at the begining, you can enqueue
   #    all the elements before creating workers.
   #
   #    At the same time, the dequeue function could be non-blocking
   #    (->dequeue_nb(COUNT) ). When no elements in queue, scripts will ends.
   #
   # 2) If you want to enqueue dynamiclly later, you do not have to enqueue
   #    here.
   #
   #    You Must use the blocking dequeue function (->dequeue(COUNT) ) to wait
   #    elements being enqueue.
   #
   #    If you want to end scripts, just enqueue an element with terminal
   #    singal that could be recognised by workers.
   #
   # For situation 1
   # $queue->enqueue(@ARGV);

    # create workers
    for ( 1 .. $WORKER_NUM ) {
        print "Create worker $_.\n";

        # worker name "worker $i" is just for testing, you could delete it.
        threads->new( \&worker, $queue, "worker_$_" );
    }

    # For Situation 2
    # enqueue dynamiclly
    run_enqueue_thread();

    # wait until all jobs being done.
    $_->join for threads->list;
}

# producer
sub run_enqueue_thread {

    sub producer {
        for (@ARGV) {
            print "Enqueue: $_.\n";
            $queue->enqueue($_);
            sleep 1;
        }
    }
    my $thread = threads->new( \&producer );

    # wait for enqueue_thread being finished.
    while (1) {
        if ( $thread->is_joinable() ) {
            $thread->join;
            last;
        }

        # check every second
        sleep 1;
    }

    print "All tasks been sended.\nSend terimnal signal.\n";

    # send terimnal signal
    for ( 1 .. $WORKER_NUM ) {
        $queue->enqueue("STOP");
    }
}

# worker could be treated as a consumer
sub worker {
    my $queue = shift;

    # $worker_name is just for testing, you could delete it.
    my $worker_name = shift;

    # Non-blocking dequeue: ->dequeue_nb(COUNT)
    # blocking dequeue    : ->dequeue(COUNT)
    # See details:
    #   http://search.cpan.org/~jdhedden/Thread-Queue-3.02/lib/Thread/Queue.pm
    #
    # This script use blocking dequeue, it will wait until queue has elements.
    while ( my $element = $queue->dequeue(1) ) {

        # recognise the terminal signal
        if ( $element eq "STOP" ) {
            print "Worker $worker_name stoped.\n";
            return 0;
        }

        # do some thing
        my $result = &DoSomething($element);
        print "Result from Worker $worker_name: $result\n";
    }
    
    return 0;
}

sub DoSomething($) {
    my $element = shift;
    sleep 2;
    return sprintf "I will do something with %s", $element;
}
