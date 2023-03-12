# Code generated by protoc-gen-twirp_ruby 1.9.0, DO NOT EDIT.
require 'twirp'
require_relative 'fftb_pb.rb'

module Fftb
  class DealerService < Twirp::Service
    service 'Dealer'
    rpc :FinishTask, FinishTaskRequest, Empty, :ruby_method => :finish_task
    rpc :QuitTask, QuitTaskRequest, Empty, :ruby_method => :quit_task
    rpc :FailTask, FailTaskRequest, Empty, :ruby_method => :fail_task
    rpc :FindFreeTask, FindFreeTaskRequest, Task, :ruby_method => :find_free_task
    rpc :Notify, NotifyRequest, Empty, :ruby_method => :notify
  end

  class DealerClient < Twirp::Client
    client_for DealerService
  end
end
