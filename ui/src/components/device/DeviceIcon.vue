<template>
  <fragment>
    <v-icon
      class="icons"
      data-test="type-icon"
    >
      {{ deviceIcon[iconName] || 'fl-tux' }}
    </v-icon>
  </fragment>
</template>

<script>

export default {
  name: 'DeviceIconComponent',

  props: {
    iconName: {
      type: String,
      required: true,
    },
  },

  data() {
    return {
      deviceIcon: {
        alpine: 'fl-alpine',
        arch: 'fl-archlinux',
        centos: 'fl-centos',
        coreos: 'fl-coreos',
        debian: 'fl-debian',
        devuan: 'fl-devuan',
        elementary: 'fl-elementary',
        fedora: 'fl-fedora',
        freebsd: 'fl-freebsd',
        gentoo: 'fl-gentoo',
        linuxmint: 'fl-linuxmint',
        mageia: 'fl-mageia',
        manjaro: 'fl-manjaro',
        mandriva: 'fl-mandriva',
        nixos: 'fl-nixos',
        opensuse: 'fl-opensuse',
        rhel: 'fl-redhat',
        sabayon: 'fl-sabayon',
        slackware: 'fl-slackware',
        ubuntu: 'fl-ubuntu',
        raspbian: 'fl-raspberry-pi',
        'ubuntu-core': 'fl-ubuntu',
        void: 'fl-void',
      },
    };
  },
};

</script>
