#!/usr/bin/env python3
"""Plot huus benchmark comparison between the main (no-cache) and feat/cache runs.

Reads benchmark_main.csv and benchmark_cache.csv from this folder and draws
INSERT / READ / DELETE latency (µs/op) vs number of keys, on log-log axes.

Usage:
    python3 plot_benchmarks.py            # writes benchmark_comparison.png
    python3 plot_benchmarks.py --show     # also opens an interactive window
"""
import argparse
import csv
import os
import sys

import matplotlib

matplotlib.use("Agg")
import matplotlib.pyplot as plt

HERE = os.path.dirname(os.path.abspath(__file__))
PHASES = [
    ("insert_us_op", "INSERT"),
    ("read_us_op", "READ"),
    ("delete_us_op", "DELETE"),
]


def load(path):
    keys, cols = [], {c: [] for c, _ in PHASES}
    with open(path, newline="") as f:
        for row in csv.DictReader(f):
            keys.append(int(row["keys"]))
            for c, _ in PHASES:
                cols[c].append(float(row[c]))
    return keys, cols


def main():
    ap = argparse.ArgumentParser()
    ap.add_argument("--show", action="store_true", help="open an interactive window")
    ap.add_argument("--out", default=os.path.join(HERE, "benchmark_comparison.png"))
    args = ap.parse_args()

    datasets = [
        ("main (no cache)", "benchmark_main.csv", "o", "#9C2A2A"),
        ("feat/cache", "benchmark_cache.csv", "s", "#305496"),
    ]
    loaded = []
    for label, fname, marker, color in datasets:
        path = os.path.join(HERE, fname)
        if not os.path.exists(path):
            sys.exit(f"missing csv: {path}")
        loaded.append((label, marker, color, *load(path)))

    fig, axes = plt.subplots(1, 3, figsize=(15, 5), sharex=True)
    for ax, (col, title) in zip(axes, PHASES):
        for label, marker, color, keys, cols in loaded:
            ax.plot(keys, cols[col], marker=marker, color=color, label=label)
        ax.set_xscale("log")
        ax.set_yscale("log")
        ax.set_title(title)
        ax.set_xlabel("keys")
        ax.set_ylabel("µs/op")
        ax.grid(True, which="both", ls=":", alpha=0.5)
        ax.legend()

    fig.suptitle("huus benchmark: main vs feat/cache (latency per op, log-log)")
    fig.tight_layout()
    fig.savefig(args.out, dpi=120)
    print("wrote", args.out)
    if args.show:
        plt.show()


if __name__ == "__main__":
    main()
