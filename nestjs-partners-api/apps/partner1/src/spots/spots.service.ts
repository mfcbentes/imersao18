import { Injectable } from '@nestjs/common';
import { CreateSpotRequest } from './request/create-spot.request';
import { UpdateSpotRequest } from './request/update-spot.request';
import { SpotStatus } from '@prisma/client';
import { PrismaService } from '@app/core/prisma/prisma.service';

@Injectable()
export class SpotsService {
  constructor(private prismaService: PrismaService) {}

  async create(createSpotRequest: CreateSpotRequest & { eventId: string }) {
    const event = await this.prismaService.event.findFirst({
      where: {
        id: createSpotRequest.eventId,
      },
    });

    if (!event) {
      throw new Error('Event not found');
    }

    return this.prismaService.spot.create({
      data: {
        ...createSpotRequest,
        status: SpotStatus.Available,
      },
    });
  }

  findAll(eventId: string) {
    return this.prismaService.spot.findMany({
      where: {
        eventId: eventId,
      },
    });
  }

  findOne(eventId: string, spotId: string) {
    return this.prismaService.spot.findFirst({
      where: {
        id: spotId,
        eventId: eventId,
      },
    });
  }

  update(
    eventId: string,
    spotId: string,
    updateSpotRequest: UpdateSpotRequest,
  ) {
    return this.prismaService.spot.update({
      where: {
        id: spotId,
        eventId: eventId,
      },
      data: updateSpotRequest,
    });
  }

  remove(eventId: string, spotId: string) {
    return this.prismaService.spot.delete({
      where: {
        id: spotId,
        eventId: eventId,
      },
    });
  }
}
