# Stat / Series / Small / Area
# Create a small stat card displaying a primary value and a series plot.
# #stat_card #series
# ---
import time

from faker import Faker

from synth import FakeCategoricalSeries
from h2o_wave import site, ui, data

page = site['/demo']

colors = '$red $pink $blue $azure $cyan $teal $mint $green $lime $yellow $amber $orange $tangerine'.split()
curves = 'linear smooth step step-after step-before'.split()
fake = Faker()
cards = []
for i in range(len(curves)):
    f = FakeCategoricalSeries()
    cat, val, pc = f.next()
    c = page.add(f'example{i}', ui.small_series_stat_card(
        box=f'1 {i + 1} 1 1',
        title=fake.cryptocurrency_name(),
        value='=${{intl qux minimum_fraction_digits=2 maximum_fraction_digits=2}}',
        data=dict(qux=val, quux=pc),
        plot_category='foo',
        plot_type='area',
        plot_value='qux',
        plot_color=colors[i],
        plot_data=data('foo qux', -15),
        plot_zero_value=0,
        plot_curve=curves[i],
    ))
    cards.append((f, c))
page.save()

while True:
    time.sleep(1)
    for f, c in cards:
        cat, val, pc = f.next()
        c.data.qux = val
        c.data.quux = pc
        c.plot_data[-1] = [cat, val]
    page.save()
