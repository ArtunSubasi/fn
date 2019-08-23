/*
 * Copyright Camunda Services GmbH and/or licensed to Camunda Services GmbH under
 * one or more contributor license agreements. See the NOTICE file distributed
 * with this work for additional information regarding copyright ownership.
 * Licensed under the Zeebe Community License 1.0. You may not use this file
 * except in compliance with the Zeebe Community License 1.0.
 */
package io.zeebe.test.util.record;

import io.zeebe.protocol.record.Record;
import io.zeebe.protocol.record.value.MessageStartEventSubscriptionRecordValue;
import java.util.stream.Stream;

public class MessageStartEventSubscriptionRecordStream
    extends ExporterRecordStream<
        MessageStartEventSubscriptionRecordValue, MessageStartEventSubscriptionRecordStream> {

  public MessageStartEventSubscriptionRecordStream(
      Stream<Record<MessageStartEventSubscriptionRecordValue>> wrappedStream) {
    super(wrappedStream);
  }

  @Override
  protected MessageStartEventSubscriptionRecordStream supply(
      Stream<Record<MessageStartEventSubscriptionRecordValue>> wrappedStream) {
    return new MessageStartEventSubscriptionRecordStream((wrappedStream));
  }

  public MessageStartEventSubscriptionRecordStream withWorkfloKey(long workflowKey) {
    return valueFilter(v -> v.getWorkflowKey() == workflowKey);
  }

  public MessageStartEventSubscriptionRecordStream withStartEventId(String startEventId) {
    return valueFilter(v -> startEventId.equals(v.getStartEventId()));
  }

  public MessageStartEventSubscriptionRecordStream withMessageName(String messageName) {
    return valueFilter(v -> messageName.equals(v.getMessageName()));
  }
}
