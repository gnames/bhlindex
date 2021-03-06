# Generated by the protocol buffer compiler.  DO NOT EDIT!
# source: protob.proto

require 'google/protobuf'

Google::Protobuf::DescriptorPool.generated_pool.build do
  add_message "protob.Version" do
    optional :value, :string, 1
  end
  add_message "protob.Void" do
  end
  add_message "protob.Item" do
    optional :id, :int32, 1
    optional :archive_id, :string, 2
    optional :path, :string, 3
    optional :lang, :string, 4
  end
  add_message "protob.ItemsOpt" do
  end
  add_message "protob.Page" do
    optional :id, :string, 1
    optional :offset, :int32, 2
    optional :text, :bytes, 3
    optional :item_id, :string, 4
    optional :item_path, :string, 5
    repeated :names, :message, 6, "protob.NameString"
  end
  add_message "protob.PagesOpt" do
    optional :with_text, :bool, 1
    repeated :item_ids, :int32, 2
  end
  add_message "protob.NameString" do
    optional :value, :string, 1
    optional :odds, :float, 2
    optional :path, :string, 3
    optional :curated, :bool, 4
    optional :edit_distance, :int32, 5
    optional :edit_distance_stem, :int32, 6
    optional :source_id, :int32, 7
    optional :match, :enum, 8, "protob.MatchType"
    optional :offset_start, :int32, 9
    optional :offset_end, :int32, 10
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
  Item = Google::Protobuf::DescriptorPool.generated_pool.lookup("protob.Item").msgclass
  ItemsOpt = Google::Protobuf::DescriptorPool.generated_pool.lookup("protob.ItemsOpt").msgclass
  Page = Google::Protobuf::DescriptorPool.generated_pool.lookup("protob.Page").msgclass
  PagesOpt = Google::Protobuf::DescriptorPool.generated_pool.lookup("protob.PagesOpt").msgclass
  NameString = Google::Protobuf::DescriptorPool.generated_pool.lookup("protob.NameString").msgclass
  MatchType = Google::Protobuf::DescriptorPool.generated_pool.lookup("protob.MatchType").enummodule
end
