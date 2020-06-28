import { Directive, OnInit } from '@angular/core';
import { Router, NavigationEnd } from '@angular/router';
import { SidebarLinkDirective } from './sidebar-link.directive';

@Directive({
  selector: '[appSidebarDropdown]'
})
export class SidebarDropdownDirective implements OnInit {
  protected navlinks: Array<SidebarLinkDirective> = [];

    public ngOnInit(): any {
      //write your code here!
    }

    constructor( private router: Router) {
    }

}
