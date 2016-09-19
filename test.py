from osmread import parse_file, Way

highway_count = 0
for entity in parse_file('europe-latest.osm.pbf'):
    if isinstance(entity, Way) and 'highway' in entity.tags:
        highway_count += 1

print("%d highways found" % highway_count)
