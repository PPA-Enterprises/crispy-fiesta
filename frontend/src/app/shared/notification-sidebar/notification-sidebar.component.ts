import { Component, OnInit, ElementRef, Renderer2, OnDestroy } from '@angular/core';
import { LayoutService } from '../services/layout.service';
import { Subscription } from 'rxjs';

@Component({
  selector: 'app-notification-sidebar',
  templateUrl: './notification-sidebar.component.html',
  styleUrls: ['./notification-sidebar.component.scss']
})
export class NotificationSidebarComponent implements OnInit, OnDestroy {

  layoutSub: Subscription;
  isOpen = false;

  ngOnInit() {

  }

  constructor(private layoutService: LayoutService) {

    this.layoutSub = layoutService.toggleNotiSidebar$.subscribe(
      open => {
        this.isOpen = open;
      });
  }

  ngOnDestroy() {
    if (this.layoutSub) {
      this.layoutSub.unsubscribe();
    }
  }

  onClose() {
    this.layoutService.toggleNotificationSidebar(false);
  }

}
