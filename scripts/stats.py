import json
import numpy as np
import matplotlib.pyplot as plt

with open('test.json') as f:
    data = json.load(f)

species_keys = ['animal', 'plant']

overtime = data['overtime']
fig, axes = plt.subplots(ncols=2, nrows=len(species_keys))
for sk in overtime.keys():
    if sk not in species_keys:
        continue

    species_dict = overtime[sk]
    for j, parkey in enumerate(species_dict.keys()):
        raw_data = species_dict[parkey]
        # assign plot index to species
        i = species_keys.index(sk)
        
        # plot data to respective subplot
        axe = axes[i,j]
        axe.plot(raw_data)
        # axe.hist(raw_data, density=1)
        axe.set_title('{}, {}'.format(sk, parkey))
        axe.set_xlim(xmin=0)



# fig, axes = plt.subplots(ncols=4, nrows=len(species_keys))
# for sk in data.keys():
#     if sk not in species_keys:
#         continue
#     species_dict = data[sk]
#     for j, parkey in enumerate(species_dict.keys()):
#         raw_data = species_dict[parkey]['all']
#         # assign plot index to species
#         i = species_keys.index(sk)
        
#         # plot data to respective subplot
#         axe = axes[i,j]
#         axe.hist(raw_data, bins=range(int(min(raw_data)), int(max(raw_data)) + 1, 1))
#         # axe.hist(raw_data, density=1)
#         axe.set_title('{}, {}'.format(sk, parkey))
#         axe.set_xlim(xmin=0)

plt.show()
