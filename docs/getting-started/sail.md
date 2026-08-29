---
title: Laravel Sail alternative
description: 'Laravel Sail gives every project its own Docker Compose stack. Lerd replaces it with one shared rootless Podman environment, and lerd import sail migrates a Sail project database and S3 files across in a single command.'
head:
  - - meta
    - name: keywords
      content: laravel sail alternative, sail without docker, migrate from laravel sail, laravel sail slow, laravel sail linux, sail vs herd, laravel local development linux, replace laravel sail
  - - script
    - type: application/ld+json
    - |
      {
        "@context": "https://schema.org",
        "@type": "FAQPage",
        "mainEntity": [
          {
            "@type": "Question",
            "name": "What is the best Laravel Sail alternative?",
            "acceptedAnswer": {
              "@type": "Answer",
              "text": "Lerd, if you work across several Laravel projects at once and do not want a Docker Compose stack per repo. It runs one shared nginx, PHP-FPM and service layer as rootless Podman containers, gives each project a .test domain and HTTPS automatically, and needs no files committed to the project. DDEV and Lando are the alternatives if you would rather stay on per-project Docker."
            }
          },
          {
            "@type": "Question",
            "name": "How do I migrate a Laravel Sail project to Lerd?",
            "acceptedAnswer": {
              "@type": "Answer",
              "text": "Run lerd sail import in the Sail project root. It starts only Sail's data services with remapped ports so nothing clashes with lerd, dumps the database, mirrors S3 and MinIO files into lerd's storage, then stops Sail. Running lerd link on a project with laravel/sail in composer.json offers the import automatically."
            }
          },
          {
            "@type": "Question",
            "name": "Do I have to delete docker-compose.yml to use Lerd?",
            "acceptedAnswer": {
              "@type": "Answer",
              "text": "No. Lerd does not touch it, so Sail keeps working. The only file lerd modifies is .env, when you run lerd env, and it saves the original as .env.before_lerd first. Run lerd env:restore to switch a project back to Sail."
            }
          },
          {
            "@type": "Question",
            "name": "Why is Laravel Sail heavy to run?",
            "acceptedAnswer": {
              "@type": "Answer",
              "text": "Each Sail project is a full Compose stack: its own PHP container, its own database, its own Redis. Five projects open means five of everything. Lerd shares one nginx, one PHP-FPM per version and one instance of each service across every site, so five running projects cost around 200 MB rather than one to two gigabytes."
            }
          },
          {
            "@type": "Question",
            "name": "Does Lerd need Docker?",
            "acceptedAnswer": {
              "@type": "Answer",
              "text": "No. Lerd runs on rootless Podman with no daemon and no root. Docker can stay installed alongside it, which is what makes the Sail import possible in the first place: lerd drives your existing Docker or Podman Compose to pull the data out."
            }
          }
        ]
      }
---

# Laravel Sail alternative

[Laravel Sail](https://laravel.com/docs/sail) is the official per-project Docker Compose setup: a `docker-compose.yml` in the repo, a container per service, and `./vendor/bin/sail` in front of every command. It works, and it is the right answer when a team wants the environment described inside the project.

The friction shows up when you have several projects open. Each one is a full stack, so five repos means five copies of PHP, five databases and five Redis instances, plus a port to remember for each. **Lerd takes the shared-infrastructure route instead**: one nginx, one PHP-FPM per version, one instance of each service, across every site, as rootless Podman containers under your own user. Every project gets a `.test` domain and a trusted certificate with nothing committed to the repo.

```bash
curl -fsSL https://lerd.sh/install.sh | bash
cd ~/code/myapp
lerd link
```

`https://myapp.test`, no port, no `sail up`.

## The one-command migration

Sail is the only tool with an automated importer in Lerd, because a Sail stack has a predictable shape. From the project root:

```bash
lerd sail import
```

It reads the Compose file, remaps any ports that would clash with Lerd's running services, starts **only** the data services (so a slow or broken app image build never blocks you), waits for the database, dumps it, imports it into Lerd's MySQL or PostgreSQL, mirrors S3 and MinIO files into Lerd's storage, then stops Sail again.

You rarely have to remember the command: running `lerd link` on a project with `laravel/sail` in its `composer.json` offers the import before setup.

The flags, the credential handling and the `.env.before_lerd` behaviour are all covered in [Importing from Laravel Sail](/usage/import-sail).

## Every Sail habit, and the Lerd equivalent

| What you did with Sail | The same thing in Lerd |
|---|---|
| `sail up -d` per project before you can work | `lerd start` once, then every project is served |
| `sail artisan`, `sail composer`, `sail npm` | `lerd artisan`, `lerd composer`, `lerd node`, from your own shell |
| `localhost:${APP_PORT}`, a different port per project | `https://myapp.test`, automatic, no ports |
| `/etc/hosts` edits when you wanted a real hostname | Automatic `.test` domains through a dnsmasq container |
| No TLS, or mkcert wired in by hand | `lerd secure`, a real mkcert certificate trusted by your system and browsers |
| Change the PHP version by switching the Sail image and rebuilding | `lerd isolate 8.4`, or let it read `composer.json`, no rebuild |
| Add a service by editing `docker-compose.yml` | `lerd service start meilisearch`, shared across every site |
| Queue and scheduler as extra Compose services | `lerd worker start queue` / `schedule`, as user services with [self-healing](/usage/worker-heal) |
| `sail logs -f` | A [log viewer](/features/logs) in the dashboard, or `lerd logs` |
| `docker-compose.yml` committed to the repo | [`.lerd.yaml`](/configuration#per-project-config-lerd-yaml), optional, a handful of lines |
| Docker Desktop, Orbstack or Colima | Rootless Podman, no daemon, no licence |

## Lerd vs Laravel Sail

|  | Lerd | Laravel Sail |
|---|---|---|
| Platforms | Linux (systemd), macOS, Windows via WSL2 (beta) | Linux, macOS, Windows |
| License | Open source (MIT) | Open source (MIT) |
| Container runtime | Rootless Podman, no daemon | Docker Desktop / Orbstack / Colima |
| Architecture | One shared nginx, PHP-FPM and service layer | A per-project Compose stack |
| PHP versions | 7.4, 8.0 to 8.5, per project, no rebuild | Per-project Sail image |
| Services (MySQL, Redis…) | One shared instance | Per project |
| `.test` domains | Automatic, zero config | Manual hosts entries, or `localhost:${APP_PORT}` |
| HTTPS | `lerd secure`, trusted mkcert certificate | Manual, or roll your own mkcert |
| RAM with 5 projects running | ~200 MB | ~1–2 GB, five full stacks |
| Requires changes to project files | No | Yes, `docker-compose.yml` committed |
| Works on legacy / client repos | Yes, just `lerd link` | Only if you can add Sail |
| Non-PHP projects | First-class via [`Containerfile.lerd`](/getting-started/containers) | Add your own container to the stack |
| Per-project service versions | No, services are shared and versioned globally | Yes, each project pins its own |
| Dashboard | [Web UI](/features/web-ui), [system tray](/features/system-tray), [terminal dashboard](/features/tui) | CLI + Docker Desktop |
| AI / MCP | Built-in [MCP server](/features/mcp) | Not built in |
| Cost | Free | Free |

**Choose Sail when:** your team already standardised on it, two projects need different major versions of the same database at the same time, or you want the infrastructure defined in the repo.

**Choose Lerd when:** you work across many projects at once, you cannot modify a client's project files, you want `.test` routing and HTTPS without wiring them, or you want the same environment on Linux and macOS.

## Running both

They coexist. Lerd uses rootless Podman with no daemon, so Docker stays installed and Sail keeps working on the projects you have not moved.

The only file Lerd touches in a project is `.env`, and only when you run `lerd env`. The original is saved as `.env.before_lerd` the first time, so a project can go back:

```bash
lerd env:restore
```

That makes the migration reversible per project, which is the sane way to move a team over: convert one repo, live on it for a week, then do the rest.

## What is genuinely different

- **Services are shared, not per project.** One MySQL, one Redis, one Mailpit across every site. That is where the memory saving comes from, and it is the one thing Sail does that Lerd does not: pinning a different database major version per project.
- **No `sail` prefix.** `lerd artisan` and `lerd composer` run in the project's PHP container but from your own shell and your own directory.
- **The environment is not in the repo.** Instead of a Compose file, you commit an optional [`.lerd.yaml`](/configuration#per-project-config-lerd-yaml) describing PHP, Node, services and workers. Teammates without Lerd are unaffected, since nothing else changed.
- **One `sudo` at install time.** Only to point the system resolver at the `.test` domains.

## Frequently asked questions

**Does the import handle S3 and MinIO?**
Yes, files are mirrored into Lerd's S3-compatible storage. Skip it with `--skip-s3` if the project does not use it.

**What if my Sail app image will not build?**
It does not matter. The import starts only the data services with `--no-deps`, so the app container is never built.

**Can I keep Sail's database credentials?**
The import auto-detects the database name and tries Sail's defaults. Override with `--sail-db-name`, `--sail-db-user` and `--sail-db-password` when your project uses custom ones.

**Does it work if I run Sail on Podman Compose?**
Yes. The command tries `docker compose` first and falls back to `podman compose`.

**Is Lerd free for commercial work?**
Yes. MIT licensed, no paid tier.

## Next steps

- [Importing from Laravel Sail](/usage/import-sail), the full flag and credential reference
- [Requirements](/getting-started/requirements) and [installation](/getting-started/installation)
- [Quick start](/getting-started/quick-start), a project served in two commands
- [Full comparison](/getting-started/comparison) against Herd, Laragon, Laradock, DDEV and Lando
