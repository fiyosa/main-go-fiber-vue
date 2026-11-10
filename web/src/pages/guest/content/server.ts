export interface ServerSpec {
  headline: string
  title: string
  summary: string
  user: string
  host: string
  os: string
  kernel: string
  uptime: string
  packages: string
  shell: string
  display: string
  terminal: string
  cpu: string
  gpu: string
  memory: string
  swap: string
  diskRoot: string
  diskLog: string
  localIp: string
  locale: string
  deployment: {
    name: string
    description: string
    details: string[]
  }
}

export const serverContent: ServerSpec = {
  headline: 'Production Host Environment',
  title: 'Server Specifications',
  summary:
    'Live telemetry and hardware specifications of the Armbian development board hosting this application, securely exposed via Cloudflare Tunnel.',
  user: 'armbian@armbian',
  host: 'Amlogic Meson GXL (S905X) P212 Development Board',
  os: 'Armbian 26.5.1 noble aarch64',
  kernel: 'Linux 5.15.195-ophub',
  uptime: '3 days, 21 hours, 12 mins',
  packages: '722 (dpkg)',
  shell: 'bash 5.2.21',
  display: 'Composite-1 (720x576, 50 Hz)',
  terminal: 'ttyd',
  cpu: 'meson-gxl (4 cores) @ 1.51 GHz',
  gpu: 'Amlogic meson-gxl-mali [Integrated]',
  memory: '732.59 MiB / 1.76 GiB (41% used)',
  swap: '174.50 MiB / 900.46 MiB (19% used)',
  diskRoot: '12.11 GiB / 28.00 GiB (43% used) - ext4',
  diskLog: '31.34 MiB / 46.84 MiB (67% used) - ext4',
  localIp: 'Hidden (LAN only)',
  locale: 'id_ID.UTF-8',
  deployment: {
    name: 'Cloudflare Tunnel',
    description:
      'Secure outbound-only tunneling connecting the local Armbian server directly to Cloudflare edge network.',
    details: [
      'No port forwarding or public IP exposure required on the local router.',
      'Encrypted HTTPS connection routed securely to the Fiber backend and Vite frontend.',
      'Edge DDoS protection, TLS termination, and global caching via Cloudflare.',
    ],
  },
}
