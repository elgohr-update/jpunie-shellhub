import Vuex from 'vuex';
import { mount, createLocalVue } from '@vue/test-utils';
import Vuetify from 'vuetify';
import TerminalDialog from '@/components/terminal/TerminalDialog';

describe('TerminalDialog', () => {
  const localVue = createLocalVue();
  const vuetify = new Vuetify();
  localVue.use(Vuex);

  let wrapper;

  const uid = 'a582b47a';
  const rules = { required: 'Required' };
  const tabs = ['Password', 'PublicKey'];
  const username = 'user';
  const passwd = 'pass';

  const privateKeys = [
    {
      name: 'shellhub',
      data: 'BBGVvbmF',
    },
    {
      name: 'shellhub',
      data: 'AbGVvbmF',
    },
  ];

  const store = new Vuex.Store({
    namespaced: true,
    state: {
      terminal: uid,
      privateKeys,
    },
    getters: {
      'modals/terminal': (state) => state.terminal,
      'privatekeys/list': (state) => state.privateKeys,
    },
    actions: {
      'modals/toggleTerminal': () => {
      },
    },
  });

  ///////
  // In this case, the rendering of the console icon is tested.
  // For this test to work, the uid in props is an empty string.
  ///////

  describe('Button', () => {
    beforeEach(() => {
      wrapper = mount(TerminalDialog, {
        store,
        localVue,
        stubs: ['fragment'],
        propsData: { uid: '', show: false, enableConnectButton: true },
        vuetify,
      });
    });

    ///////
    // Component Rendering
    //////

    it('Is a Vue instance', () => {
      document.body.setAttribute('data-app', true);
      expect(wrapper).toBeTruthy();
    });
    it('Renders the component', () => {
      expect(wrapper.html()).toMatchSnapshot();
    });

    ///////
    // Data and Props checking
    //////

    it('Receive data in props', () => {
      expect(wrapper.vm.uid).toEqual('');
      expect(wrapper.vm.show).toEqual(false);
    });
    it('Compare data with default value', () => {
      wrapper.setData({ rules });

      expect(wrapper.vm.username).toEqual('');
      expect(wrapper.vm.passwd).toEqual('');
      expect(wrapper.vm.showLoginForm).toEqual(true);
      expect(wrapper.vm.valid).toEqual(true);
      expect(wrapper.vm.rules).toEqual(rules);
      expect(wrapper.vm.tabs).toEqual(tabs);
    });
    it('Receive data in computed', () => {
      expect(wrapper.vm.showTerminal).toEqual(false);
      expect(wrapper.vm.show).toEqual(false);
      expect(wrapper.vm.getListPrivateKeys).toEqual(privateKeys);
    });
    it('Check the watch action', async () => {
      wrapper.setData({ username, passwd });

      await wrapper.vm.$options.watch.show.call(wrapper.vm, false);

      expect(wrapper.vm.username).toEqual(username);
      expect(wrapper.vm.passwd).toEqual(passwd);
    });

    //////
    // HTML validation
    //////

    it('Renders the template with data', () => {
      expect(wrapper.find('[data-test="connect-btn"]').exists()).toBe(true);
      expect(wrapper.find('[data-test="console-icon"]').exists()).toBe(false);
      expect(wrapper.find('[data-test="console-item"]').exists()).toBe(false);
      expect(wrapper.find('[data-test="terminal-dialog"]').exists()).toBe(false);
    });
  });

  describe('Icon', () => {
    beforeEach(() => {
      wrapper = mount(TerminalDialog, {
        store,
        localVue,
        stubs: ['fragment'],
        propsData: { uid: '', show: false },
        vuetify,
      });
    });

    ///////
    // Component Rendering
    //////

    it('Is a Vue instance', () => {
      document.body.setAttribute('data-app', true);
      expect(wrapper).toBeTruthy();
    });
    it('Renders the component', () => {
      expect(wrapper.html()).toMatchSnapshot();
    });

    ///////
    // Data and Props checking
    //////

    it('Receive data in props', () => {
      expect(wrapper.vm.uid).toEqual('');
      expect(wrapper.vm.show).toEqual(false);
    });
    it('Compare data with default value', () => {
      wrapper.setData({ rules });

      expect(wrapper.vm.username).toEqual('');
      expect(wrapper.vm.passwd).toEqual('');
      expect(wrapper.vm.showLoginForm).toEqual(true);
      expect(wrapper.vm.valid).toEqual(true);
      expect(wrapper.vm.rules).toEqual(rules);
      expect(wrapper.vm.tabs).toEqual(tabs);
    });
    it('Receive data in computed', () => {
      expect(wrapper.vm.showTerminal).toEqual(false);
      expect(wrapper.vm.show).toEqual(false);
      expect(wrapper.vm.getListPrivateKeys).toEqual(privateKeys);
    });
    it('Check the watch action', async () => {
      wrapper.setData({ username, passwd });

      await wrapper.vm.$options.watch.show.call(wrapper.vm, false);

      expect(wrapper.vm.username).toEqual(username);
      expect(wrapper.vm.passwd).toEqual(passwd);
    });

    //////
    // HTML validation
    //////

    it('Renders the template with data', () => {
      expect(wrapper.find('[data-test="connect-btn"]').exists()).toBe(false);
      expect(wrapper.find('[data-test="console-icon"]').exists()).toBe(true);
      expect(wrapper.find('[data-test="console-item"]').exists()).toBe(true);
      expect(wrapper.find('[data-test="terminal-dialog"]').exists()).toBe(false);
    });
  });

  ///////
  // In this case the dialog is opened
  ///////

  describe('Dialog opened', () => {
    beforeEach(() => {
      wrapper = mount(TerminalDialog, {
        store,
        localVue,
        stubs: ['fragment'],
        propsData: { uid, show: true, enableConnectButton: true },
        vuetify,
      });
    });

    ///////
    // Component Rendering
    //////

    it('Is a Vue instance', () => {
      document.body.setAttribute('data-app', true);
      expect(wrapper).toBeTruthy();
    });
    it('Renders the component', () => {
      expect(wrapper.html()).toMatchSnapshot();
    });

    ///////
    // Data and Props checking
    //////

    it('Receive data in props', () => {
      expect(wrapper.vm.uid).toEqual(uid);
      expect(wrapper.vm.show).toEqual(true);
    });
    it('Compare data with default value', () => {
      wrapper.setData({ rules });

      expect(wrapper.vm.username).toEqual('');
      expect(wrapper.vm.passwd).toEqual('');
      expect(wrapper.vm.showLoginForm).toEqual(true);
      expect(wrapper.vm.valid).toEqual(true);
      expect(wrapper.vm.rules).toEqual(rules);
      expect(wrapper.vm.tabs).toEqual(tabs);
    });
    it('Receive data in computed', () => {
      expect(wrapper.vm.showTerminal).toEqual(true);
      expect(wrapper.vm.show).toEqual(true);
      expect(wrapper.vm.getListPrivateKeys).toEqual(privateKeys);
    });
    it('Check the watch action', async () => {
      wrapper.setData({ username, passwd });
      await wrapper.vm.$options.watch.show.call(wrapper.vm, false);

      expect(wrapper.vm.username).toEqual(username);
      expect(wrapper.vm.passwd).toEqual(passwd);
    });

    //////
    // HTML validation
    //////

    it('Renders the template with data', async () => {
      expect(wrapper.find('[data-test="connect-btn"]').exists()).toBe(true);
      expect(wrapper.find('[data-test="console-icon"]').exists()).toBe(false);
      expect(wrapper.find('[data-test="console-item"]').exists()).toBe(false);
      // Validate if the table was created
      tabs.forEach((item) => {
        expect(wrapper.find(`[data-test="${item}-tab"]`).exists()).toEqual(true);
      });

      // When tab is password
      expect(wrapper.find('[data-test="terminal-dialog"]').exists()).toBe(true);
      expect(wrapper.find('[data-test="username-field"]').exists()).toBe(true);
      expect(wrapper.find('[data-test="passwd-field"]').exists()).toBe(true);
      expect(wrapper.find('[data-test="connect-btn"]').exists()).toBe(true);

      // When tab is publicKey
      wrapper.find('[data-test="PublicKey-tab"]').trigger('click');
      await wrapper.vm.$nextTick();
      await wrapper.vm.$nextTick();

      expect(wrapper.find('[data-test="username2-field"]').exists()).toBe(true);
      expect(wrapper.find('[data-test="privatekeys-select"]').exists()).toBe(true);
      expect(wrapper.find('[data-test="connect2-btn"]').exists()).toBe(true);
    });
  });
});
