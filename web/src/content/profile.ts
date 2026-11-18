// Portfolio content for the guest `/about` page.
// Every string is sourced from the ATS CV; do not add facts that are not there.

export interface ProfileProject {
  name: string
  description: string
  bullets: string[]
}

export interface ProfileRole {
  role: string
  company: string
  period: string
  blurb: string
  bullets: string[]
}

export interface ProfileSkillGroup {
  group: string
  items: string[]
}

export interface ProfileEducation {
  school: string
  degree: string
  gpa: string
  period: string
}

export interface Profile {
  name: string
  headline: string
  location: string
  email: string
  summary: string
  roles: ProfileRole[]
  projects: ProfileProject[]
  skills: ProfileSkillGroup[]
  education: ProfileEducation
}

export const profile: Profile = {
  name: 'Fiter Yogie Saputra',
  headline: 'Fullstack Developer',
  location: 'Bandung, Indonesia',
  email: 'fiterys@gmail.com',
  summary:
    'Fullstack Developer with 5+ years of experience designing and building business web applications and dashboards end to end. Comfortable adapting across different frontend and backend stacks depending on project needs and team requirements, with a strong ability to pick up new tools quickly. Focused on delivering reliable, maintainable systems from database design through deployment, while keeping close attention to performance and user experience on business-facing products throughout.',
  roles: [
    {
      role: 'Full Stack Software Engineer',
      company: 'PT Sigma Cipta Caraka',
      period: 'Jan 2022 - Present',
      blurb:
        'An IT services and consulting company delivering end to end technology solutions for enterprise and general business clients across multiple industries.',
      bullets: [
        'Developed and maintained internal web applications and dashboards, working across both frontend and backend layers.',
        'Built user interfaces and backend services to support internal business workflows and reporting needs.',
        'Collaborated with cross functional teams to gather requirements and deliver features for internal stakeholders.',
        'Contributed to improving application reliability, performance, and code quality across ongoing projects.',
      ],
    },
    {
      role: 'Backend Engineer',
      company: 'PT Daya Desa Asasta',
      period: 'Mar 2021 - Oct 2021',
      blurb:
        'An IT company developing Internet of Things (IoT) based technology solutions to support internal business operations and workflows.',
      bullets: [
        'Designed and developed backend services and APIs supporting internal web applications and dashboards.',
        'Built and maintained database schemas and business logic to support internal reporting and operational workflows.',
        'Collaborated with frontend developers to integrate backend services and ensure smooth data flow across systems.',
        'Contributed to improving API performance, code quality, and system reliability across ongoing projects.',
      ],
    },
  ],
  projects: [
    {
      name: 'KYPAS',
      description:
        'An internal web application built for Telkom Indonesia employees to manage job position data and organizational role structures, while also handling creation and management of employee VPN access.',
      bullets: [
        'Built fullstack features for managing employee job position data and organizational role structures across Telkom Indonesia.',
        'Developed functionality to create and manage VPN access for employees based on their roles and permissions.',
        'Designed and implemented both frontend interfaces and backend APIs to support employee and access management workflows.',
        'Managed user data at scale, ensuring accurate role assignments and consistent access records across the organization.',
      ],
    },
    {
      name: 'SSO',
      description:
        'A centralized Single Sign-On system built for Telkom Indonesia employees, enabling a single login to seamlessly access multiple internal applications without requiring repeated authentication across different systems.',
      bullets: [
        'Designed and developed the fullstack SSO system, covering both frontend login interfaces and backend authentication services.',
        'Implemented centralized authentication logic to allow seamless access across multiple connected internal applications.',
        'Built and maintained APIs to handle user session management and token-based authentication.',
        'Collaborated with teams across different applications to integrate SSO login, reducing repeated authentication steps for employees.',
      ],
    },
  ],
  skills: [
    { group: 'Languages', items: ['JavaScript', 'TypeScript', 'Python', 'Go', 'SQL'] },
    { group: 'Frontend', items: ['React.js', 'Next.js', 'Vue', 'SvelteKit', 'Tailwind CSS'] },
    { group: 'Backend', items: ['Laravel', 'Gin', 'ASP.NET Core', 'NestJS', 'Express.js', 'REST APIs'] },
    {
      group: 'Databases & Infrastructure',
      items: ['PostgreSQL', 'MySQL', 'Supabase', 'JWT Authentication', 'API integrations'],
    },
    {
      group: 'DevOps',
      items: [
        'Docker',
        'CI/CD',
        'Linux server administration and deployment',
        'OKD (OpenShift Kubernetes Distribution)',
        'Harbor',
      ],
    },
    { group: 'Product & Integrations', items: ['Real-time communication', 'Dashboard development'] },
  ],
  education: {
    school: 'State Polytechnic of Bandung',
    degree: 'Diploma IV Electronics Engineering',
    gpa: 'GPA 3.60 / 4.00',
    period: 'Aug 2016 - Nov 2020',
  },
}
