export default {
  /* eslint-disable no-param-reassign */
  install(Vue) {
    const env = {
      isEnterprise: (window.env || process.env).VUE_APP_SHELLHUB_ENTERPRISE === 'true',
      isCloud: (window.env || process.env).VUE_APP_SHELLHUB_CLOUD === 'true',
      stripePublishableKey: (window.env || process.env).VUE_APP_SHELLHUB_STRIPE_PUBLISHABLE_KEY,
      billingEnable: (window.env || process.env).VUE_APP_SHELLHUB_BILLING === 'true',
    };

    Vue.env = env;
    Vue.prototype.$env = env;
  },
  /* eslint-enable no-param-reassign */
};
