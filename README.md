# pgfathom

Sound the depth of a legacy PostgreSQL schema.

`pgfathom` discovers the relationships your database has but never declared — and proves them against the data instead of guessing from column names.

> **Status: pre-release.** Nothing is implemented yet. The design is specified in
> [`docs/PGFATHOM.md`](docs/PGFATHOM.md). Recovery-rate benchmarks will be
> published here once the MVP runs against the reference corpus.

---

## The problem

Old PostgreSQL databases carry more structure than they declare. `pedido.cliente_id`
points at `cliente.id` in every single row, but no constraint says so — the ORM of the
day never created it, or someone dropped constraints for a bulk load and never put them
back.

Two things follow. Nobody can read the model, because `\d` shows no relationships at all.
And because the constraint was never there, nothing ever stopped orphan rows from getting
in — so they probably already did, years ago, silently.

There is a nastier variant: a foreign key can be declared and still guarantee nothing,
if it was created `NOT VALID` and never validated. It shows up in `\d`, it draws an arrow
in your ERD tool, and it never checked a single pre-existing row.

## What it does

`pgfathom` reads the system catalog, mines join predicates out of view and function
definitions, filters candidates against planner statistics, and then validates what
survives with an anti-join against the real data.

Every relationship it reports comes with the evidence behind it: containment by row and
by distinct value, orphan counts, and cardinality. Relationships fall into three buckets
that need three different responses — a forgotten but intact foreign key, a real
relationship with broken integrity, or a name coincidence.

The second bucket is the point of the tool. It is a data bug that has been in production
for years and nobody knows.

Output is a terminal report, a versioned JSON model, and reviewable `.sql` artifacts.

## Guarantees

**Read-only.** `pgfathom` never issues a statement that modifies the database under
analysis. It emits `.sql` files for you to review and run yourself.

**Your data never leaves memory.** The tool reads values to compare keys. What comes out
of it are counts, ratios, and object names — never a value from your tables, in any
output, log, or JSON field. The target use case includes public-sector databases holding
national ID numbers, health records, and taxpayer data.

**No claim without evidence.** Every inferred relationship carries a verdict and the
metric backing it. Sampled runs can never report a confirmed relationship — orphan rows
cluster on disk, which is exactly the pattern page-level sampling is worst at finding.

**Silence is never reported as a clean bill of health.** Tables skipped for missing
privileges, candidates that timed out, schemas not covered — all of it appears in the
coverage block on every run.

## Prior art

Containment is known in the data-profiling literature as an *inclusion dependency*, and
there is a mature body of work on discovering them — SPIDER, BINDER, MIND, all implemented
in [Metanome](https://hpi.de/naumann/projects/repeatability/data-profiling/metanome-ind-algorithms.html).
Commercial GUI modelers such as Hackolade infer relationships from metadata.
[Azimutt](https://github.com/azimuttapp/azimutt) flags `_id` columns without declared
relations as part of its schema analysis.

`pgfathom` is not new science. It is the part nobody packaged: a native PostgreSQL CLI
that validates inference against the actual data, mines evidence from the catalog itself,
speaks the naming conventions of legacy schemas in languages other than English, and hands
you DDL you can review.

## License

TBD.
