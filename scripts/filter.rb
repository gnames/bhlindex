# frozen_string_literal: true

require 'csv'

names = {}

puts "Filtering names.csv to names_filtered.csv\n\n"

# copy headers of names first
r = File.open('names.csv')
w = File.open('names_filtered.csv', 'w:utf-8')
w.write(r.readline)
r.close
w.close

# apply a filter to names
w = CSV.open('names_filtered.csv', 'a:utf-8')
r = CSV.open('names.csv', headers: true)

count = 0
saved = 0
r.each do |l|
  count += 1
  puts "Traversed #{count} names, saved #{saved}" if (count % 1_000_000).zero?
  next unless %w[Curated AutoCurated].include?(l['Curation']) || l['OddsLog10'].to_f > 6.0

  w << l
  names[l['NameID']] = true
  saved += 1
end

r.close
w.close

puts "\nFiltering occurrences.csv to occurrences_filtered.csv\n\n"

r = File.open('occurrences.csv')
w = File.open('occurrences_filtered.csv', 'w:utf-8')
saved = 0
r.each_with_index do |l, i|
  if i.zero?
    w.write(l)
    next
  end
  id = l[0..35]
  puts "Traversed #{i} occurrences, saved #{saved}" if (i % 10_000_000).zero?
  if names[id]
    saved += 1
    w.write(l)
  end
end

r.close
w.close
