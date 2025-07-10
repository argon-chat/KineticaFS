# KineticaFS
KineticaFS is a region-aware file lifecycle manager for S3-based storage. It tracks file usage across regions, migrates hot files closer to clients, and safely removes unused ones using reference-counting and GC.


## ✨ Features

- 📦 **S3-compatible**: Works with any S3-compatible storage backend (AWS, MinIO, Wasabi, etc.)
- 🌍 **Region-aware**: Detects file access patterns and migrates "hot" files closer to clients
- 🧠 **Smart pointers**: Tracks file references in your system to prevent premature deletion
- ♻️ **Garbage collection**: Removes unreferenced or expired files safely and automatically


## 📈 Roadmap

- [ ] File reference tracking API (`CreateRef`, `DeleteRef`, `ListRefs`)
- [ ] File upload and migration logic
- [ ] Per-region heatmap tracking
- [ ] GC for unreferenced files
- [ ] Basic observability (logs, metrics)
- [ ] Public and expiring file links
- [ ] Optional TTL per reference
- [ ] Support for batch import/export
- [ ] Multi-tenant support
- [ ] NATS hook for event-driven GC
- [ ] Custom metadata indexing
- [ ] Integration with Prometheus / Grafana
- [ ] WASM hooks for file filters / pre-upload logic
- [ ] Admin panel / metrics endpoint

## 📜 License

KineticaFS is licensed under the **GNU Affero General Public License v3.0 or later (AGPL-3.0-or-later)**.  
This ensures that all improvements and deployments based on this code must remain open source.

See [`LICENSE`](./LICENSE) for the full text.
