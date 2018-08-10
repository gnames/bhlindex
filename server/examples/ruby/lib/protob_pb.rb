# Generated by the protocol buffer compiler.  DO NOT EDIT!
# source: protob.proto

require 'google/protobuf'

Google::Protobuf::DescriptorPool.generated_pool.build do
  add_message "protob.Version" do
    optional :value, :string, 1
  end
  add_message "protob.Void" do
  end
  add_message "protob.Title" do
    optional :id, :string, 1
    optional :path, :string, 2
    repeated :pages, :message, 3, "protob.Page"
  end
  add_message "protob.Page" do
    optional :id, :string, 1
    repeated :names, :message, 2, "protob.NameString"
  end
  add_message "protob.NameString" do
    optional :value, :string, 1
    optional :odds, :float, 2
    optional :path, :string, 3
    optional :curated, :bool, 4
    optional :edit_distance, :int32, 5
    optional :edit_distance_stem, :int32, 6
    optional :match, :enum, 7, "protob.MatchType"
  end
  add_enum "protob.MatchType" do
    value :NONE, 0
    value :EXACT, 1
    value :CANONICAL_EXACT, 2
    value :CANONICAL_FUZZY, 3
    value :PARTIAL_EXACT, 4
    value :PARTIAL_FUZZY, 5
  end
end

module Protob
  Version = Google::Protobuf::DescriptorPool.generated_pool.lookup("protob.Version").msgclass
  Void = Google::Protobuf::DescriptorPool.generated_pool.lookup("protob.Void").msgclass
  Title = Google::Protobuf::DescriptorPool.generated_pool.lookup("protob.Title").msgclass
  Page = Google::Protobuf::DescriptorPool.generated_pool.lookup("protob.Page").msgclass
  NameString = Google::Protobuf::DescriptorPool.generated_pool.lookup("protob.NameString").msgclass
  MatchType = Google::Protobuf::DescriptorPool.generated_pool.lookup("protob.MatchType").enummodule
end