<script setup lang="ts">
  import { BrutalBadge, BrutalCard } from '../../components/guest/ui'
  import { profile } from '../../content/profile'

  // Anchors share the BrutalButton shape without nesting a button inside a link.
  // The contact anchor must not uppercase, so the email stays lowercase.
  const primaryLinkClasses =
    'brutal-press inline-flex min-h-12 items-center justify-center gap-2 border-[3px] border-fg bg-vivid-yellow px-6 font-mono text-sm font-bold tracking-widest text-black shadow-brutal-sm'

  const secondaryLinkClasses =
    'brutal-press inline-flex min-h-12 items-center justify-center gap-2 border-[3px] border-fg bg-surface px-6 font-mono text-sm font-bold tracking-widest text-fg uppercase shadow-brutal-sm'

  const projectAccents = ['sky', 'soft-violet'] as const
  const skillTones = ['sky', 'mint', 'soft-violet', 'vivid-yellow', 'hot-red'] as const
</script>

<template>
  <div class="flex flex-col gap-14">
    <!-- Hero -->
    <section class="grid items-start gap-8 lg:grid-cols-[1.4fr_0.6fr]" aria-labelledby="profile-title">
      <div class="flex flex-col items-start gap-4">
        <BrutalBadge tone="hot-red" :rotation="-2">{{ profile.headline }}</BrutalBadge>

        <h1 id="profile-title" class="text-4xl leading-none font-bold tracking-tight sm:text-5xl">
          {{ profile.name }}
        </h1>

        <p class="max-w-2xl text-base leading-relaxed">
          Business web applications and dashboards, designed and built
          <span class="bg-vivid-yellow px-1 text-black">end to end</span>.
        </p>

        <div class="flex flex-wrap items-center gap-3">
          <a href="#projects" :class="secondaryLinkClasses">See projects</a>
        </div>
      </div>

      <BrutalCard accent="vivid-yellow" tone="muted" class="w-full rotate-1">
        <template #header>Based in</template>

        <p class="text-xl font-bold">{{ profile.location }}</p>
        <p class="mt-2 font-mono text-xs break-all">{{ profile.email }}</p>
      </BrutalCard>
    </section>

    <!-- Summary -->
    <section id="summary" class="scroll-mt-24" aria-labelledby="summary-title">
      <h2 id="summary-title" class="text-2xl font-bold tracking-tight sm:text-3xl">Summary</h2>
      <p class="mt-4 max-w-3xl text-base leading-relaxed">{{ profile.summary }}</p>
    </section>

    <!-- Experience -->
    <section id="experience" class="scroll-mt-24" aria-labelledby="experience-title">
      <h2 id="experience-title" class="text-2xl font-bold tracking-tight sm:text-3xl">Experience</h2>

      <ol class="mt-6 flex flex-col gap-6 border-l-[3px] border-fg pl-6 sm:pl-8">
        <li v-for="role in profile.roles" :key="role.company" class="relative">
          <span
            aria-hidden="true"
            class="absolute top-7 -left-[calc(1.5rem+11px)] size-4 border-[3px] border-fg bg-vivid-yellow sm:-left-[calc(2rem+11px)]"
          ></span>

          <BrutalCard tone="surface" class="w-full">
            <template #header>{{ role.company }}</template>

            <div class="flex flex-wrap items-center gap-3">
              <h3 class="text-xl font-bold">{{ role.role }}</h3>
              <BrutalBadge tone="mint">{{ role.period }}</BrutalBadge>
            </div>

            <p class="mt-3 text-sm leading-relaxed">{{ role.blurb }}</p>

            <ul class="mt-4 flex flex-col gap-2">
              <li v-for="bullet in role.bullets" :key="bullet" class="flex gap-3 text-sm leading-relaxed">
                <span aria-hidden="true" class="mt-2 size-2 shrink-0 border-2 border-fg bg-hot-red"></span>
                <span>{{ bullet }}</span>
              </li>
            </ul>
          </BrutalCard>
        </li>
      </ol>
    </section>

    <!-- Projects -->
    <section id="projects" class="scroll-mt-24" aria-labelledby="projects-title">
      <h2 id="projects-title" class="text-2xl font-bold tracking-tight sm:text-3xl">Projects</h2>

      <div class="mt-6 grid gap-6 md:grid-cols-2">
        <BrutalCard
          v-for="(project, index) in profile.projects"
          :key="project.name"
          :accent="projectAccents[index % 2]"
        >
          <template #header>{{ project.name }}</template>

          <p class="text-sm leading-relaxed">{{ project.description }}</p>

          <ul class="mt-4 flex flex-col gap-2">
            <li v-for="bullet in project.bullets" :key="bullet" class="flex gap-3 text-sm leading-relaxed">
              <span aria-hidden="true" class="mt-2 size-2 shrink-0 border-2 border-fg bg-vivid-yellow"></span>
              <span>{{ bullet }}</span>
            </li>
          </ul>
        </BrutalCard>
      </div>
    </section>

    <!-- Skills -->
    <section id="skills" class="scroll-mt-24" aria-labelledby="skills-title">
      <h2 id="skills-title" class="text-2xl font-bold tracking-tight sm:text-3xl">Skills</h2>

      <div class="mt-6 grid gap-6 sm:grid-cols-2 lg:grid-cols-3">
        <BrutalCard v-for="(group, index) in profile.skills" :key="group.group" tone="surface">
          <template #header>{{ group.group }}</template>

          <ul class="flex flex-wrap gap-2">
            <li v-for="item in group.items" :key="item">
              <BrutalBadge :tone="skillTones[index % skillTones.length]">{{ item }}</BrutalBadge>
            </li>
          </ul>
        </BrutalCard>
      </div>
    </section>

    <!-- Education -->
    <section id="education" class="scroll-mt-24" aria-labelledby="education-title">
      <h2 id="education-title" class="text-2xl font-bold tracking-tight sm:text-3xl">Education</h2>

      <BrutalCard accent="soft-violet" class="mt-6">
        <div class="flex flex-wrap items-center gap-3">
          <h3 class="text-xl font-bold">{{ profile.education.school }}</h3>
          <BrutalBadge tone="surface">{{ profile.education.period }}</BrutalBadge>
        </div>

        <p class="mt-3 text-sm leading-relaxed">{{ profile.education.degree }}</p>
        <p class="mt-1 font-mono text-xs tracking-widest uppercase">{{ profile.education.gpa }}</p>
      </BrutalCard>
    </section>

    <!-- Contact -->
    <section id="contact" class="scroll-mt-24" aria-labelledby="contact-title">
      <BrutalCard accent="hot-red" tone="muted">
        <template #header>Contact</template>

        <div class="flex flex-col gap-5">
          <h2 id="contact-title" class="text-2xl font-bold tracking-tight sm:text-3xl">Let's build something solid.</h2>

          <div class="flex flex-wrap items-center gap-3">
            <a :href="`mailto:${profile.email}`" :class="primaryLinkClasses">{{ profile.email }}</a>
            <BrutalBadge tone="surface">{{ profile.location }}</BrutalBadge>
          </div>
        </div>
      </BrutalCard>
    </section>
  </div>
</template>
