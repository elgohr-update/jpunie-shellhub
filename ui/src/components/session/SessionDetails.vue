<template>
  <fragment>
    <div class="d-flex pa-0 align-center">
      <h1 v-if="hide">
        Session Details
      </h1>
      <v-spacer />
      <v-spacer />
    </div>

    <v-card
      v-if="session"
      class="mt-2"
    >
      <v-toolbar
        flat
        color="transparent"
      >
        <v-toolbar-title v-if="session.device">
          <v-icon
            v-if="session.active"
            color="success"
          >
            check_circle
          </v-icon>
          <v-tooltip
            v-else
            bottom
          >
            <template #activator="{ on }">
              <v-icon v-on="on">
                check_circle
              </v-icon>
            </template>
            <span>active {{ session.last_seen | lastSeen }}</span>
          </v-tooltip>
          {{ session.device.name }}
        </v-toolbar-title>

        <v-spacer />

        <v-menu
          v-if="session.recorded"
          ref="menu"
          offset-y
        >
          <template #activator="{ on, attrs }">
            <v-chip
              color="transparent"
              v-on="on"
            >
              <v-icon
                small
                class="icons"
                v-bind="attrs"
                v-on="on"
              >
                mdi-dots-horizontal
              </v-icon>
            </v-chip>
          </template>

          <v-card>
            <v-tooltip
              bottom
              :disabled="hasAuthorizationPlay"
            >
              <template #activator="{ on, attrs }">
                <div
                  v-bind="attrs"
                  v-on="on"
                >
                  <v-list-item
                    v-if="session.recorded && isEnterprise"
                    :disabled="!hasAuthorizationPlay"
                    @click.stop="openDialog('sessionPlayDialog')"
                  >
                    <SessionPlay
                      :uid="session.uid"
                      :recorded="session.authenticated && session.recorded"
                      :show.sync="sessionPlayDialog"
                      data-test="sessionPlay-component"
                    />
                  </v-list-item>
                </div>
              </template>

              <span>
                You don't have this kind of authorization.
              </span>
            </v-tooltip>

            <v-tooltip
              bottom
              :disabled="hasAuthorizationClose"
            >
              <template #activator="{ on, attrs }">
                <div
                  v-bind="attrs"
                  v-on="on"
                >
                  <v-list-item
                    v-if="session.active"
                    :disabled="!hasAuthorizationClose"
                    @click.stop="openDialog('sessionCloseDialog')"
                  >
                    <SessionClose
                      :uid="session.uid"
                      :device="session.device_uid"
                      :show.sync="sessionCloseDialog"
                      data-test="sessionClose-component"
                      @update="refresh"
                    />
                  </v-list-item>
                </div>
              </template>

              <span>
                You don't have this kind of authorization.
              </span>
            </v-tooltip>

            <v-tooltip
              bottom
              :disabled="hasAuthorizationRemoveRecord"
            >
              <template #activator="{ on, attrs }">
                <div
                  v-bind="attrs"
                  v-on="on"
                >
                  <v-list-item
                    v-if="session.recorded"
                    :disabled="!hasAuthorizationRemoveRecord"
                    @click.stop="openDialog('sessionDeleteRecord')"
                  >
                    <SessionDeleteRecord
                      :uid="session.uid"
                      :show.sync="sessionDeleteRecord"
                      data-test="sessionDeleteRecord-component"
                      @update="refresh"
                    />
                  </v-list-item>
                </div>
              </template>

              <span>
                You don't have this kind of authorization.
              </span>
            </v-tooltip>
          </v-card>
        </v-menu>
      </v-toolbar>

      <v-divider />

      <v-card-text>
        <div class="mt-2">
          <div class="overline">
            Uid
          </div>
          <div
            data-test="sessionUid-field"
          >
            {{ session.uid }}
          </div>
        </div>

        <div class="mt-2">
          <div class="overline">
            User
          </div>
          <div
            data-test="sessionUser-field"
          >
            {{ session.username }}
          </div>
        </div>

        <div class="mt-2">
          <div
            class="overline"
          >
            Authenticated
          </div>
          <v-tooltip bottom>
            <template
              v-if="session"
              #activator="{ on }"
            >
              <v-icon
                v-if="session.authenticated"
                :color="session.active ? 'success' : ''"
                size=""
                v-on="on"
              >
                mdi-shield-check
              </v-icon>
              <v-icon
                v-else
                color="error"
                size=""
                v-on="on"
              >
                mdi-shield-alert
              </v-icon>
            </template>
            <span v-if="session.authenticated">User has been authenticated</span>
            <span v-else>User has not been authenticated</span>
          </v-tooltip>
        </div>

        <div class="mt-2">
          <div class="overline">
            Ip Address
          </div>
          <code
            data-test="sessionIpAddress-field"
          >
            {{ session.ip_address }}
          </code>
        </div>

        <div class="mt-2">
          <div class="overline">
            Started
          </div>
          <div
            data-test="sessionStartedAt-field"
          >
            {{ session.started_at | formatDate }}
          </div>
        </div>

        <div class="mt-2">
          <div class="overline">
            Last Seen
          </div>
          <div
            data-test="sessionLastSeen-field"
          >
            {{ session.last_seen | formatDate }}
          </div>
        </div>
      </v-card-text>
    </v-card>

    <div class="text-center">
      <v-dialog
        v-model="dialog"
        width="500"
        persistent
      >
        <v-card>
          <v-card-title class="headline primary">
            Session ID error
          </v-card-title>
          <v-card-text class="mt-4 mb-3 pb-1">
            You tried to access a non-existing session ID.
          </v-card-text>
          <v-card-actions>
            <v-spacer />
            <v-btn
              color="primary"
              text
              @click="redirect"
            >
              Go back to sessions
            </v-btn>
          </v-card-actions>
        </v-card>
      </v-dialog>
    </div>
  </fragment>
</template>

<script>

import SessionPlay from '@/components/session/SessionPlay';
import SessionClose from '@/components/session/SessionClose';
import SessionDeleteRecord from '@/components/session/SessionDeleteRecord';
import { formatDate, lastSeen } from '@/components/filter/date';
import hasPermission from '@/components/filter/permission';

export default {
  name: 'SessionDetailsComponent',

  components: {
    SessionPlay,
    SessionClose,
    SessionDeleteRecord,
  },

  filters: { formatDate, lastSeen, hasPermission },

  data() {
    return {
      uid: '',
      session: null,
      dialog: false,
      sessionPlayDialog: false,
      sessionCloseDialog: false,
      sessionDeleteRecord: false,
      hide: true,
      playAction: 'play',
      closeAction: 'close',
      removeRecordAction: 'removeRecord',
    };
  },

  computed: {
    isEnterprise() {
      return this.$env.isEnterprise;
    },

    hasAuthorizationPlay() {
      const role = this.$store.getters['auth/role'];
      if (role !== '') {
        return hasPermission(
          this.$authorizer.role[role],
          this.$actions.session[this.playAction],
        );
      }

      return false;
    },

    hasAuthorizationClose() {
      const role = this.$store.getters['auth/role'];
      if (role !== '') {
        return hasPermission(
          this.$authorizer.role[role],
          this.$actions.session[this.closeAction],
        );
      }

      return false;
    },

    hasAuthorizationRemoveRecord() {
      const role = this.$store.getters['auth/role'];
      if (role !== '') {
        return hasPermission(
          this.$authorizer.role[role],
          this.$actions.session[this.removeRecordAction],
        );
      }

      return false;
    },
  },

  async created() {
    this.uid = this.$route.params.id;
    try {
      await this.$store.dispatch('sessions/get', this.uid);
      this.session = this.$store.getters['sessions/get'];
    } catch (error) {
      this.hide = false;
      this.dialog = true;
      this.$store.dispatch('snackbar/showSnackbarErrorLoading', this.$errors.snackbar.sessionDetails);
    }
  },

  methods: {
    redirect() {
      this.dialog = false;
      this.$router.push('/sessions');
    },

    async refresh() {
      try {
        await this.$store.dispatch('sessions/get', this.uid);
        this.session = this.$store.getters['sessions/get'];
      } catch {
        this.$store.dispatch('snackbar/showSnackbarErrorLoading', this.$errors.snackbar.sessionDetails);
      }
    },

    openDialog(action) {
      this[action] = !this[action];
      this.closeMenu();
    },

    closeMenu() {
      this.$refs.menu.isActive = false;
    },
  },
};

</script>
