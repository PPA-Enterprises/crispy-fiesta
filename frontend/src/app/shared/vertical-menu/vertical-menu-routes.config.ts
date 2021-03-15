import { RouteInfo } from './vertical-menu.metadata';

//Sidebar menu Routes and data
export const ROUTES: RouteInfo[] = [
  {
    path: '/home', title: 'Home', icon: 'ft-home', class: 'dropdown nav-item', isExternalLink: false, submenu: []
  },
  {
    path: '/jobs', title: 'Jobs', icon: 'ft-truck', class: 'dropdown nav-item', isExternalLink: false, submenu: []
  },
  {
    path: '/clients', title: 'Clients', icon: 'ft-user', class: 'dropdown nav-item', isExternalLink: false, submenu: []
  },
  {
    path: '/calendar', title: 'Calendar', icon: 'ft-user', class: 'dropdown nav-item', isExternalLink: false, submenu: []
  },
//   {
//     path: '', title: 'Menu Levels', icon: 'ft-align-left', class: 'has-sub', badge: '3', badgeClass: 'badge badge-pill badge-danger float-right mr-1 mt-1', isExternalLink: false,
//     submenu: [
//         { path: '/YOUR-ROUTE-PATH', title: 'Second Level', icon: 'ft-arrow-right submenu-icon', class: '', badge: '', badgeClass: '', isExternalLink: false, submenu: [] },
//         {
//             path: '', title: 'Second Level Child', icon: 'ft-arrow-right submenu-icon', class: 'has-sub', badge: '', badgeClass: '', isExternalLink: false,
//             submenu: [
//                 { path: '/YOUR-ROUTE-PATH', title: 'Third Level 1.1', icon: 'ft-arrow-right submenu-icon', class: '', badge: '', badgeClass: '', isExternalLink: false, submenu: [] },
//                 { path: '/YOUR-ROUTE-PATH', title: 'Third Level 1.2', icon: 'ft-arrow-right submenu-icon', class: '', badge: '', badgeClass: '', isExternalLink: false, submenu: [] },
//             ]
//         },
//     ]
// },
];
