
/* === CLASS GENERATORS === */
const validOptions = {
      boolean: [true, false],
      booleanProps: [
        'loading',
        'disabled',
        'error',
        'transparent',
        'inverted',
        'fluid',
        'readonly',
        'checked',
        'slider',
        'toggle',
        'required',
        'multiple',
        'inline',
        'success',
        'warning',
        'styled',
        'multiple',
        'simple',
        'animated',
        'animating',
        'looping',
        'scrolling',
        'compact',
        'attached',
        'animateOnLoad'
      ],

      size: ['mini', 'tiny', 'small', 'medium', 'large', 'big', 'huge', 'massive'],
      color: ['red', 'orange', 'yellow', 'olive', 'green', 'teal', 'blue', 'violet', 'purple', 'pink', 'brown', 'grey', 'black'],

      buttonType: ['basic', 'animated', 'labeled', 'icon'],
      inputType: ['text', 'password', 'email', 'number', 'tel', 'url', 'date', 'datetime-local', 'time', 'month', 'week', 'search', 'file'],
      tableType: ['striped', 'celled', 'padded', 'compact', 'basic', 'structured'],
      
      menuType: ['vertical', 'horizontal', 'tabular', 'pagination', 'secondary', 'pointing', 'compact', 'text', 'icon', 'labeled icon'],
      menuPosition: ['top', 'bottom', 'left', 'right'],

      transition: ['scale', 'fade', 'fade up', 'fade down', 'fade left', 'fade right',
                 'horizontal flip', 'vertical flip', 'drop', 'fly left', 'fly right',
                 'fly up', 'fly down', 'swing left', 'swing right', 'swing up',
                 'swing down', 'browse', 'browse right', 'slide down', 'slide up',
                 'slide left', 'slide right'],
      transitionDirection: ['in', 'out'],

      columnSize: [
        'one', 'two', 'three', 'four',
        'five', 'six', 'seven', 'eight',
        'nine', 'ten', 'eleven', 'twelve',
        'thirteen', 'fourteen', 'fifteen', 'sixteen'
      ],

      sidebarDirection: ['top', 'bottom', 'left', 'right'],
      sidebarVariation: ['overlay', 'push', 'scale down', 'uncover', 'slide along', 'slide out'],

      popupPosition: ['top left', 'top center', 'top right', 
                'bottom left', 'bottom center', 'bottom right',
                'left center', 'right center'],

      popupVariation: ['basic', 'wide', 'very wide', 'flowing', 'mini', 'tiny', 'small', 'large'],

      modalSize: ['mini', 'tiny', 'small', 'medium', 'large', 'fullscreen'],
      modalPosition: ['top', 'bottom', 'left', 'right'],

      dropdownDirection: ['upward', 'downward', 'left', 'right'],
      dropdownAction: ['hide', 'show', 'toggle'],
 }

const ClassGenerators = {
  // Button Generator
  button: (props) => {
    const valid = {
      size: ['mini', 'tiny', 'small', 'medium', 'large', 'big', 'huge', 'massive'],
      color: ['red', 'orange', 'yellow', 'olive', 'green', 'teal', 'blue', 'violet', 'purple', 'pink', 'brown', 'grey', 'black'],
      type: ['basic', 'animated', 'labeled', 'icon'],
    };
    // const booleanVals = ['fluid', 'animated', 'loading', 'disabled'];
    return ['ui', 'button']
      .concat(valid.size.includes(props.size) ? props.size : '')
      .concat(valid.color.includes(props.color) ? props.color : '')
      .concat(valid.type.includes(props.type) ? props.type : '')
      .concat(props.fluid ? 'fluid' : '')
      .concat(props.animated ? 'animated' : '')
      .concat(props.loading ? 'loading' : '')
      .concat(props.disabled ? 'disabled' : '')
      .filter(Boolean).join(' ');
  },

  // Input Generator
  input: (props) => {
    const valid = {
      size: ['mini', 'tiny', 'small', 'medium', 'large', 'big', 'huge', 'massive'],
      color: ['red', 'orange', 'yellow', 'olive', 'green', 'teal', 'blue', 'violet', 'purple', 'pink', 'brown', 'grey', 'black'],
      type: ['text', 'password', 'email', 'number', 'tel', 'url', 'date', 'datetime-local', 'time', 'month', 'week', 'search', 'file']
    };
    return ['ui', 'input', 'field']
      .concat(valid.size.includes(props.size) ? props.size : '')
      .concat(valid.color.includes(props.color) ? props.color : '')
      .concat(props.type && valid.type.includes(props.type) ? props.type : '')
      .concat(props.loading ? 'loading' : '')
      .concat(props.disabled ? 'disabled' : '')
      .concat(props.error ? 'error' : '')
      .concat(props.transparent ? 'transparent' : '')
      .concat(props.inverted ? 'inverted' : '')
      .concat(props.fluid ? 'fluid' : '')
      .filter(Boolean).join(' ');
  },

  // Radio Button Generator
  radio: (props) => {
    return ['ui', 'radio', 'checkbox']
      .concat(props.checked ? 'checked' : '')
      .concat(props.disabled ? 'disabled' : '')
      .concat(props.readonly ? 'readonly' : '')
      .filter(Boolean).join(' ');
  },

  // Select Dropdown Generator
  select: (props) => {
    const valid = {
      size: ['mini', 'tiny', 'small', 'medium', 'large', 'big', 'huge', 'massive'],
      searchSize: ['one', 'two', 'three', 'four', 'five', 'six', 'seven', 'eight', 'nine', 'ten']
    };
    return ['ui', 'dropdown', 'selection']
      .concat(valid.size.includes(props.size) ? props.size : '')
      .concat(props.multiple ? 'multiple' : '')
      .concat(props.search ? 'search' : '')
      .concat(props.searchSize && valid.searchSize.includes(props.searchSize) ? props.searchSize : '')
      .concat(props.disabled ? 'disabled' : '')
      .concat(props.loading ? 'loading' : '')
      .concat(props.error ? 'error' : '')
      .filter(Boolean).join(' ');
  },

  // Checkbox Generator
  checkbox: (props) => {
    return ['ui', 'checkbox']
      .concat(props.checked ? 'checked' : '')
      .concat(props.disabled ? 'disabled' : '')
      .concat(props.readonly ? 'readonly' : '')
      .concat(props.toggle ? 'toggle' : '')
      .concat(props.slider ? 'slider' : '')
      .filter(Boolean).join(' ');
  },

  // Form Field Generator
  field: (props) => {
    return ['ui', 'field']
      .concat(props.inline ? 'inline' : '')
      .concat(props.disabled ? 'disabled' : '')
      .concat(props.error ? 'error' : '')
      .concat(props.required ? 'required' : '')
      .filter(Boolean).join(' ');
  },

  // Textarea Generator
  textarea: (props) => {
    return ['ui', 'textarea']
      .concat(props.disabled ? 'disabled' : '')
      .concat(props.error ? 'error' : '')
      .concat(props.rows ? `rows-${props.rows}` : '')
      .concat(props.resize === false ? 'no-resize' : '')
      .filter(Boolean).join(' ');
  },

  // Form Generator
  form: (props) => {
    return ['ui', 'form']
      .concat(props.loading ? 'loading' : '')
      .concat(props.error ? 'error' : '')
      .concat(props.success ? 'success' : '')
      .concat(props.warning ? 'warning' : '')
      .concat(props.inverted ? 'inverted' : '')
      .filter(Boolean).join(' ');
  },

  // Card Generator
  card: (props) => {
    return ['ui', 'card']
      .concat(props.color || '')
      .concat(props.raised ? 'raised' : '')
      .concat(props.centered ? 'centered' : '')
      .filter(Boolean).join(' ');
  },
  /* === ACCORDION === */
  accordion: (props) => {
    return ['ui', 'accordion']
      .concat(props.styled ? 'styled' : '')
      .concat(props.fluid ? 'fluid' : '')
      .concat(props.inverted ? 'inverted' : '')
      .filter(Boolean).join(' ');
  },

  /* === DROPDOWN === */
  dropdown: (props) => {
    const valid = {
      direction: ['upward', 'downward', 'left', 'right'],
      action: ['hide', 'show', 'toggle'],
      color: ['red', 'orange', 'yellow', 'olive', 'green', 'teal', 'blue', 'violet', 'purple', 'pink', 'brown', 'grey', 'black']
    };
    return ['ui', 'dropdown']
      .concat(props.selection ? 'selection' : '')
      .concat(props.search ? 'search' : '')
      .concat(props.multiple ? 'multiple' : '')
      .concat(props.color && valid.color.includes(props.color) ? props.color : '')
      .concat(props.direction && valid.direction.includes(props.direction) ? props.direction : '')
      .concat(props.simple ? 'simple' : '')
      .concat(props.fluid ? 'fluid' : '')
      .concat(props.floating ? 'floating' : '')
      .concat(props.scrolling ? 'scrolling' : '')
      .concat(props.compact ? 'compact' : '')
      .filter(Boolean).join(' ');
  },

  /* === MODAL === */
  modal: (props) => {
    const valid = {
      size: ['mini', 'tiny', 'small', 'medium', 'large', 'fullscreen'],
      position: ['top', 'bottom', 'left', 'right']
    };
    return ['ui', 'modal']
      .concat(props.size && valid.size.includes(props.size) ? props.size : '')
      .concat(props.position && valid.position.includes(props.position) ? `${props.position} aligned` : '')
      .concat(props.inverted ? 'inverted' : '')
      .concat(props.basic ? 'basic' : '')
      .concat(props.scrolling ? 'scrolling' : '')
      .concat(props.dimmer ? `dimmer ${props.dimmer}` : '') // e.g. "blurring", "inverted"
      .filter(Boolean).join(' ');
  },

  /* === POPUP === */
  popup: (props) => {
    const valid = {
      position: ['top left', 'top center', 'top right', 
                'bottom left', 'bottom center', 'bottom right',
                'left center', 'right center'],
      variation: ['basic', 'wide', 'very wide', 'flowing', 'mini', 'tiny', 'small', 'large']
    };
    return ['ui', 'popup']
      .concat(props.position && valid.position.includes(props.position) ? props.position : '')
      .concat(props.variation && valid.variation.includes(props.variation) ? props.variation : '')
      .concat(props.inverted ? 'inverted' : '')
      .concat(props.fluid ? 'fluid' : '')
      .filter(Boolean).join(' ');
  },

  /* === PROGRESS === */
  progress: (props) => {
    const valid = {
      color: ['red', 'orange', 'yellow', 'olive', 'green', 'teal', 'blue', 'violet', 'purple', 'pink', 'brown', 'grey', 'black'],
      size: ['tiny', 'small', 'medium', 'large', 'big']
    };
    return ['ui', 'progress']
      .concat(props.color && valid.color.includes(props.color) ? props.color : '')
      .concat(props.size && valid.size.includes(props.size) ? props.size : '')
      .concat(props.active ? 'active' : '')
      .concat(props.success ? 'success' : '')
      .concat(props.warning ? 'warning' : '')
      .concat(props.error ? 'error' : '')
      .concat(props.disabled ? 'disabled' : '')
      .concat(props.indeterminate ? 'indeterminate' : '')
      .concat(props.attached ? `${props.attached} attached` : '')
      .filter(Boolean).join(' ');
  },

  /* === RATING === */
  rating: (props) => {
    return ['ui', 'rating']
      .concat(props.disabled ? 'disabled' : '')
      .concat(props.readonly ? 'readonly' : '')
      .concat(props.clearable ? 'clearable' : '')
      .concat(props.icon ? 'icon' : '')
      .concat(props.star ? 'star' : '')
      .concat(props.heart ? 'heart' : '')
      .filter(Boolean).join(' ');
  },

  /* === SIDEBAR === */
  sidebar: (props) => {
    const valid = {
      direction: ['top', 'bottom', 'left', 'right'],
      variation: ['overlay', 'push', 'scale down', 'uncover', 'slide along', 'slide out']
    };
    return ['ui', 'sidebar']
      .concat(props.direction && valid.direction.includes(props.direction) ? props.direction : '')
      .concat(props.variation && valid.variation.includes(props.variation) ? props.variation : '')
      .concat(props.inverted ? 'inverted' : '')
      .concat(props.visible ? 'visible' : '')
      .filter(Boolean).join(' ');
  },

  /* === TAB === */
  tab: (props) => {
    return ['ui', 'tab']
      .concat(props.loading ? 'loading' : '')
      .concat(props.active ? 'active' : '')
      .concat(props.attached ? `${props.attached} attached` : '')
      .filter(Boolean).join(' ');
  },

  /* === TRANSITION === */
  transition: (props) => {
    const valid = {
      animateOnLoad: [true, false],
      animation: ['scale', 'fade', 'fade up', 'fade down', 'fade left', 'fade right',
                 'horizontal flip', 'vertical flip', 'drop', 'fly left', 'fly right',
                 'fly up', 'fly down', 'swing left', 'swing right', 'swing up',
                 'swing down', 'browse', 'browse right', 'slide down', 'slide up',
                 'slide left', 'slide right'],
      direction: ['in', 'out'],
      duration: [100, 200, 300, 400, 500, 600, 700, 800, 900, 1000],
      looping: [true, false],
    };
    return ['ui', 'transition']
      .concat(props.animation && valid.animation.includes(props.animation) ? props.animation : '')
      .concat(props.direction && valid.direction.includes(props.direction) ? props.direction : '')
      .concat(props.animateOnLoad ? 'animating' : '')
      .concat(props.visible ? 'visible' : '')
      .concat(props.disabled ? 'disabled' : '')
      .concat(props.looping ? 'looping' : '')
      .filter(Boolean).join(' ');
  },

  /* === MENU === */
  menu: (props) => {
    const valid = {
      color: ['red', 'orange', 'yellow', 'olive', 'green', 'teal', 'blue', 'violet', 'purple', 'pink', 'brown', 'grey', 'black'],
      type: ['vertical', 'horizontal', 'tabular', 'pagination', 'secondary', 'pointing', 'compact', 'text', 'icon', 'labeled icon'],
      position: ['top', 'bottom', 'left', 'right']
    };
    return ['ui', 'menu']
      .concat(props.color && valid.color.includes(props.color) ? props.color : '')
      .concat(props.type && valid.type.includes(props.type) ? props.type : '')
      .concat(props.position && valid.position.includes(props.position) ? `${props.position} fixed` : '')
      .concat(props.borderless ? 'borderless' : '')
      .concat(props.fluid ? 'fluid' : '')
      .concat(props.inverted ? 'inverted' : '')
      .concat(props.attached ? `${props.attached} attached` : '') // 'top', 'bottom'
      .concat(props.stackable ? 'stackable' : '')
      .filter(Boolean).join(' ');
  },

  /* === MESSAGE === */
  message: (props) => {
    const valid = {
      color: ['red', 'orange', 'yellow', 'olive', 'green', 'teal', 'blue', 'violet', 'purple', 'pink', 'brown', 'grey', 'black'],
      size: ['mini', 'tiny', 'small', 'medium', 'large', 'big', 'huge', 'massive']
    };
    return ['ui', 'message']
      .concat(props.color && valid.color.includes(props.color) ? props.color : '')
      .concat(props.size && valid.size.includes(props.size) ? props.size : '')
      .concat(props.floating ? 'floating' : '')
      .concat(props.compact ? 'compact' : '')
      .concat(props.attached ? `${props.attached} attached` : '') // 'top', 'bottom'
      .concat(props.visible ? 'visible' : '')
      .concat(props.hidden ? 'hidden' : '')
      .concat(props.icon ? 'icon' : '')
      .filter(Boolean).join(' ');
  },

  /* === TABLE === */
  table: (props) => {
    const valid = {
      color: ['red', 'orange', 'yellow', 'olive', 'green', 'teal', 'blue', 'violet', 'purple', 'pink', 'brown', 'grey', 'black'],
      size: ['mini', 'tiny', 'small', 'medium', 'large', 'big', 'huge', 'massive'],
      type: ['striped', 'celled', 'padded', 'compact', 'basic', 'structured']
    };
    return ['ui', 'table']
      .concat(props.color && valid.color.includes(props.color) ? props.color : '')
      .concat(props.size && valid.size.includes(props.size) ? props.size : '')
      .concat(props.type && valid.type.includes(props.type) ? props.type : '')
      .concat(props.collapsing ? 'collapsing' : '')
      .concat(props.selectable ? 'selectable' : '')
      .concat(props.sortable ? 'sortable' : '')
      .concat(props.inverted ? 'inverted' : '')
      .concat(props.singleLine ? 'single line' : '')
      .concat(props.fixed ? 'fixed' : '')
      .concat(props.unstackable ? 'unstackable' : '')
      .filter(Boolean).join(' ');
  },
};

/* === COMPOSITION MONAD === */
const ClassComposer = {
  create() {
    let classes = [];
    return {
      add(generator, props) {
        classes.push(generator(props));
        return this;
      },
      when(condition, ...additionalClasses) {
        if (condition) classes.push(...additionalClasses);
        return this;
      },
      build() {
        return classes.join(' ').replace(/\s+/g, ' ').trim();
      }
    };
  }
};

/* === PETITE-VUE COMPONENTS === */
const SemanticUI = {
  // Button
  Button: (props) => ({
    state: {...props},
    get classes() {
      return ClassComposer.create()
        .add(ClassGenerators.button, props)
        .when(this.loading, 'loading')
        .build();
    },
    toggle() {
      this.loading = !this.loading;
    }
  }),

  // Input
  Input: (props) => ({
    state: {...props},
    data: {
      value: null,
    },
    get classes() {
      return ClassComposer.create()
        .add(ClassGenerators.input, props)
        .build();
    },
    validate(fn) {
      return fn(this.data.value);
    }    
  }),

  Card: (props) => ({
    state: {...props},
    get classes() {
      return ClassComposer.create()
        .add(ClassGenerators.card, props)
        .when(this.highlighted, 'raised', 'blue')
        .build();
    },
    toggle() {
      this.highlighted != this.highlighted;
    }
  }),

  Radio: (props) => ({
    state: {...props},
    get classes() {
      return ClassComposer.create()
        .add(ClassGenerators.radio, props)
        .build();
    },
  }),


  Select: (props) => ({
    state: {...props},
    get classes() {
      return ClassComposer.create()
        .add(ClassGenerators.radio, props)
        .build();
    },
  }),

  Textarea: (props) => ({
    state: {...props},
    get classes() {
      return ClassComposer.create()
        .add(ClassGenerators.radio, props)
        .build();
    },
  }),


  Form: (props) => ({
    state: {...props},
    get classes() {
      return ClassComposer.create()
        .add(ClassGenerators.radio, props)
        .build();
    },
  }),



  Field: (props) => ({
    state: {...props},
    get classes() {
      return ClassComposer.create()
        .add(ClassGenerators.field, props)
        .build();
    },
  }),

  Accordion: (props) => ({
    state: {...props},
    get classes() {
      return ClassComposer.create()
      .add(ClassGenerators.accordion, props)
      .build();
    },
  }),

  Dropdown: (props) => ({
    state: {...props},
    get classes() {
      return ClassComposer.create()
      .add(ClassGenerators.dropdown, props)
      .build();
    },
  }),

  Modal: (props) => ({
    state: {...props},
    get classes() {
      return ClassComposer.create()
      .add(ClassGenerators.modal, props)
      .build();
    },
  }),

  Popup: (props) => ({
    state: {...props},
    get classes() {
      return ClassComposer.create()
      .add(ClassGenerators.popup, props)
      .build();
    },
  }),

  Progress: (props) => ({
    state: {...props},
    get classes() {
      return ClassComposer.create()
      .add(ClassGenerators.progress, props)
      .build();
    },
  }),

  Rating: (props) => ({
    state: {...props},
    get classes() {
      return ClassComposer.create()
      .add(ClassGenerators.rating, props)
      .build();
    },
  }),

  Sidebar: (props) => ({
    state: {...props},
    get classes() {
      return ClassComposer.create()
      .add(ClassGenerators.sidebar, props)
      .build();
    },
  }),

  Tab: (props) => ({
    state: {...props},
    get classes() {
      return ClassComposer.create()
      .add(ClassGenerators.tab, props)
      .build();
    },
  }),

  Transition: (props) => ({
    state: {...props},
    get classes() {
      return ClassComposer.create()
      .add(ClassGenerators.transitionn, props)
      .build();
    },
  }),

  Menu: (props) => ({
    state: {...props},
    get classes() {
      return ClassComposer.create()
      .add(ClassGenerators.menu, props)
      .build();
    },
  }),

  Message: (props) => ({
    state: {...props},
    get classes() {
      return ClassComposer.create()
      .add(ClassGenerators.message, props)
      .build();
    },
  }),

  Table: (props) => ({
    state: {...props},
    get classes() {
      return ClassComposer.create()
      .add(ClassGenerators.table, props)
      .build();
    },
  }),


};



// Usage Examples:
// HTML Usage:

// ```html
// <div v-scope="sutton({ color: 'blue', size: 'large' })">
//   <button :class="classes" @click="toggle">
//     {{ count }} ({{ doubleCount }})
//   </button>
// </div>

// <div v-scope="SmartCard({ color: 'teal', raised: true })">
//   <div :class="classes">Card Content</div>
// </div>
// ```

// Direct Class Generation:

// ```javascript
// const buttonClass = initSemanticApp().semantic('button', { 
//   color: 'green', 
//   loading: true 
// });
// // "ui green button loading"
// ```











// // semantic-classes.js
// const SemanticClasses = {
//   classProps: {
//     size: ['mini', 'tiny', 'small', 'medium', 'large', 'big', 'huge', 'massive'],
//     color: [
//       'red', 'orange', 'yellow', 'olive', 'green', 'teal',
//       'blue', 'violet', 'purple', 'pink', 'brown', 'grey', 'black'
//     ],
//     elevation: ['raised', 'stacked', 'piled'],
//     floated: ['left', 'right'],
//     aligned: ['left', 'center', 'right', 'justified'],
//     attached: ['top', 'bottom'],
//     inverted: [true, false],
//     padded: [true, false],
//     compact: [true, false],
//     circular: [true, false],
//     loading: [true, false],
//     disabled: [true, false]
//   },

//   generateClasses(props) {
//     return Object.entries(props)
//       .filter(([prop, value]) => 
//         this.classProps[prop] && 
//         value !== undefined &&
//         value !== false &&
//         this.classProps[prop].includes(value)
//       )
//       .map(([prop, value]) => {
//         if (value === true) return prop;
//         if (typeof value === 'string') return `${value} ${prop}`;
//         return '';
//       })
//       .filter(Boolean)
//       .join(' ');
//   }
// };


// const ButtonClasses = {
//   classProps: {
//     size: ['mini', 'tiny', 'small', 'medium', 'large', 'big', 'huge', 'massive'],
//     color: ['red', 'orange', 'blue', .../* all colors */],
//     basic: [true, false],
//     inverted: [true, false],
//     animated: ['fade', 'vertical', false],
//     loading: [true, false],
//     disabled: [true, false],
//     icon: [true, 'left', 'right'],
//     labeled: [true, 'left', 'right']
//   },
//   generate(props) {
//     let base = 'ui button';
//     return base + ' ' + Object.entries(props)
//       .filter(([key, val]) => this.classProps[key]?.includes(val))
//       .map(([key, val]) => val === true ? key : `${val} ${key}`)
//       .join(' ');
//   }
// };

// const FormClasses = {
//   classProps: {
//     size: ['small', 'large'],
//     inverted: [true, false],
//     loading: [true, false],
//     equalWidth: ['fields', true, false],
//     grouped: [true, 'unstackable'],
//     unstackable: [true, false]
//   },
//   generate(props) {
//     let base = 'ui form';
//     // Special handling for field types
//     if (props.fieldType) base += ` ${props.fieldType} field`;
//     return base + ' ' + /* ...same filter/map as above */;
//   }
// };

// const CardClasses = {
//   classProps: {
//     raised: [true, false],
//     centered: [true, false],
//     fluid: [true, false],
//     link: [true, false],
//     color: ['red', 'blue', ...],
//     horizontal: [true, false]
//   },
//   generate(props) {
//     let base = 'ui card';
//     if (props.header) base += ' header';
//     return base + ' ' + /* ...filter/map */;
//   }
// };

// const MenuClasses = {
//   classProps: {
//     vertical: [true, false],
//     tabular: [true, 'right'],
//     pointing: [true, false],
//     secondary: [true, false],
//     compact: [true, false],
//     fixed: ['top', 'bottom', 'left', 'right']
//   },
//   generate(props) {
//     let base = 'ui menu';
//     if (props.icons) base += ' icon';
//     return base + ' ' + /* ... */;
//   }
// };

// const MessageClasses = {
//   classProps: {
//     type: ['info', 'success', 'warning', 'error'],
//     size: ['mini', 'tiny', 'small', 'large', 'huge'],
//     floating: [true, false],
//     compact: [true, false]
//   },
//   generate(props) {
//     let base = 'ui message';
//     if (props.type) base += ` ${props.type}`;
//     return base + ' ' + /* ... */;
//   }
// };

// const TransitionClasses = {
//   transitions: {
//     type: [
//       'scale', 'fade', 'fade up', 'fade down', 'fade left', 'fade right',
//       'horizontal flip', 'vertical flip', 'drop', 'fly left', 'fly right',
//       'fly up', 'fly down', 'swing left', 'swing right', 'swing up',
//       'swing down', 'browse', 'browse right', 'slide down', 'slide up',
//       'slide left', 'slide right'
//     ],
//     direction: ['in', 'out'],
//     duration: [100, 200, 300, 400, 500, 600, 700, 800, 900, 1000],
//     loop: [true, false],
//     animateOnLoad: [true, false]
//   },

//   generate(config) {
//     const classes = ['transition'];
    
//     // Add transition type (required)
//     if (this.classProps.type.includes(config.type)) {
//       classes.push(config.type);
//     } else {
//       classes.push('fade'); // default
//     }

//     // Add direction
//     if (config.direction && this.classProps.direction.includes(config.direction)) {
//       classes.push(config.direction);
//     }

//     // Add duration class (if specified)
//     if (config.duration && this.classProps.duration.includes(config.duration)) {
//       classes.push(`duration-${config.duration}`);
//     }

//     // Additional flags
//     if (config.loop) classes.push('loop');
//     if (config.animateOnLoad) classes.push('animating');

//     return classes.join(' ');
//   },

//   // Helper to generate jQuery initialization code
//   jqueryInit(elementSelector, options = {}) {
//     return `$('${elementSelector}').transition(${JSON.stringify(options)})`;
//   }
// };


// const ClassComposer = {
//   // Start with empty classes
//   of() {
//     return {
//       classes: [],
//       // Add classes from any generator
//       add(generator, config) {
//         this.classes.push(generator.generate(config));
//         return this;
//       },
//       // Apply modifiers (like 'inverted', 'disabled')
//       modifier(name, active = true) {
//         if (active) this.classes.push(name);
//         return this;
//       },
//       // Combine with conditional classes
//       when(condition, className) {
//         if (condition) this.classes.push(className);
//         return this;
//       },
//       // Final string output
//       build() {
//         return this.classes.join(' ').replace(/\s+/g, ' ').trim();
//       }
//     };
//   }
// };

// // Example usage:
// const buttonProps = {
//   size: 'large',
//   color: 'blue',
//   inverted: true,
//   loading: true
// };

// const buttonClasses = SemanticClasses.generateClasses(buttonProps);
// // Returns: "large blue inverted loading"
// Usage in Components:
// javascript
// // In your component file
// function Button(props) {
//   const classes = SemanticClasses.generateClasses(props);
  
//   return {
//     classes, // Expose to template
//     ...props,
//     // other component logic
//   };
// }

// In your HTML/Petite-Vue template:
// <button :class="classes" @click="handleClick">Click me</button>



// 1. Create component
// function MyButton(props) {
//   return {
//     classes: ButtonClasses.generate(props),
//     ...props
//   };
// }

// 2. Use in template
// <button v-scope="MyButton({ color: 'blue', loading: true })" 
//         :class="classes">
//   Submit
// </button>

// const combinedClasses = [
//   LayoutClasses.generate(gridProps),
//   ButtonClasses.generate(buttonProps)
// ].join(' ');

// Example Usage:
// const transitionConfig = {
//   type: 'fade up',
//   direction: 'in',
//   duration: 500,
//   animateOnLoad: true
// };

// const transitionClass = TransitionClasses.generate(transitionConfig);
// Returns: "transition fade up in duration-500 animating"

// // For jQuery initialization:
// TransitionClasses.jqueryInit('.my-element', {
//   animation: 'fade up',
//   duration: 500
// });

// function AnimatedCard(props) {
//   return {
//     transitionClass: TransitionClasses.generate(props.transition),
//     init() {
//       // Initialize with jQuery
//       $(this.$el).transition({
//         animation: props.transition.type,
//         duration: props.transition.duration || 300
//       });
//     }
//   };
// }

// <div v-scope="AnimatedCard({
//   transition: {
//     type: 'fade up',
//     duration: 500
//   }
// })" 
// :class="transitionClass">
//   <!-- Content -->
// </div>

// // Usage Example:
// const buttonClasses = ClassComposer
//   .of()
//   .add(ButtonClasses, { 
//     color: 'blue', 
//     size: 'large' 
//   })
//   .modifier('inverted', true)
//   .when(isLoading, 'loading')
//   .build();

// // Returns: "ui large blue inverted button loading"

// const cardClasses = ClassComposer
//   .of()
//   .add(CardClasses, {
//     raised: true,
//     color: 'teal'
//   })
//   .add(TransitionClasses, {
//     type: 'fade up'
//   })
//   .build();


// function SmartButton(props) {
//   return {
//     classes: ClassComposer.of()
//       .add(ButtonClasses, props)
//       .when(props.isDisabled, 'disabled')
//       .when(props.isLoading, 'loading')
//       .build(),
//     // ... component logic
//   };
// }


// function MyButton(props) {
//   return {
//     // Reactive class composition
//     get classes() {
//       return ClassComposer.of()
//         .add(ButtonClasses, props)
//         .build();
//     },
//     // Initial setup
//     $mounted() {
//       console.log("Mounted");
//     }
//   };
// }


// function MyButton() {
//   return {
//     externalState: false,
//     classes: '',
    
//     // Manual update handler
//     $update() {
//       this.classes = ClassComposer.of()
//         .add(ButtonClasses, { 
//           color: this.externalState ? 'blue' : 'grey' 
//         })
//         .build();
//     },
    
//     // Example external trigger
//     toggleState() {
//       this.externalState = !this.externalState;
//       this.$update(); // Manually trigger
//     }
//   };
// }

// function MyButton() {
//   return {
//     externalState: false,
    
//     // Reactive getter (better approach)
//     get classes() {
//       return ClassComposer.of()
//         .add(ButtonClasses, {
//           color: this.externalState ? 'blue' : 'grey'
//         })
//         .build();
//     },
    
//     toggleState() {
//       this.externalState = !this.externalState;
//       // No $update needed - Petite-Vue detects the change
//     }
//   };
// }

// <!-- Example usage -->
// <div v-scope="MyButton()" 
//      :class="classes"
//      @click="toggleState">
//   Click me
// </div>

// function SmartButton() {
//   return {
//     size: 'medium',
//     color: 'blue',
//     loading: false,
    
//     // Computed classes
//     get classes() {
//       return ClassComposer.of()
//         .add(ButtonClasses, {
//           size: this.size,
//           color: this.color
//         })
//         .when(this.loading, 'loading')
//         .build();
//     },
    
//     toggleLoading() {
//       this.loading = !this.loading;
//     }
//   };
// }

// <div v-scope="SmartButton()">
//   <button :class="classes" @click="toggleLoading">
//     Click me ({{ classes }})
//   </button>
// </div>














// Here's a single-file solution with all Semantic UI generators and Petite-Vue integration:

// ```javascript
// semantic-helpers.js
